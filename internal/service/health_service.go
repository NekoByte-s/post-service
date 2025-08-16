package service

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"postService/internal/database"
	"postService/internal/models"
)

type HealthService interface {
	GetHealth() *models.HealthResponse
	GetLiveness() *models.LivenessResponse
	GetReadiness() *models.ReadinessResponse
	CheckComponent(name string) *models.ComponentHealth
}

type healthService struct {
	db        *database.Database
	startTime time.Time
	version   string
}

func NewHealthService(db *database.Database, version string) HealthService {
	return &healthService{
		db:        db,
		startTime: time.Now(),
		version:   version,
	}
}

func (s *healthService) GetHealth() *models.HealthResponse {
	timestamp := time.Now()
	uptime := timestamp.Sub(s.startTime)
	
	// Check all components
	components := []models.ComponentHealth{
		*s.checkDatabase(),
		*s.checkMemory(),
		*s.checkGoroutines(),
	}
	
	// Calculate summary
	summary := s.calculateSummary(components)
	
	// Determine overall status
	overallStatus := s.determineOverallStatus(components)
	
	return &models.HealthResponse{
		Status:     overallStatus,
		Timestamp:  timestamp,
		Version:    s.version,
		Uptime:     int64(uptime.Seconds()),
		Components: components,
		Summary:    summary,
	}
}

func (s *healthService) GetLiveness() *models.LivenessResponse {
	// Liveness is simple - if we can respond, we're alive
	return &models.LivenessResponse{
		Status:    models.HealthStatusHealthy,
		Timestamp: time.Now(),
		Message:   "Service is alive and responding",
	}
}

func (s *healthService) GetReadiness() *models.ReadinessResponse {
	timestamp := time.Now()
	
	// Check critical components for readiness
	components := []models.ComponentHealth{
		*s.checkDatabase(),
	}
	
	// Determine readiness status
	status := models.HealthStatusHealthy
	message := "Service is ready to accept requests"
	
	for _, comp := range components {
		if comp.Status == models.HealthStatusUnhealthy {
			status = models.HealthStatusUnhealthy
			message = "Service is not ready - critical components unhealthy"
			break
		} else if comp.Status == models.HealthStatusDegraded && status == models.HealthStatusHealthy {
			status = models.HealthStatusDegraded
			message = "Service is partially ready - some components degraded"
		}
	}
	
	return &models.ReadinessResponse{
		Status:     status,
		Timestamp:  timestamp,
		Message:    message,
		Components: components,
	}
}

func (s *healthService) CheckComponent(name string) *models.ComponentHealth {
	switch name {
	case "database":
		return s.checkDatabase()
	case "memory":
		return s.checkMemory()
	case "goroutines":
		return s.checkGoroutines()
	default:
		return &models.ComponentHealth{
			Name:    name,
			Status:  models.HealthStatusUnhealthy,
			Message: "Unknown component",
			Error:   "Component not found",
		}
	}
}

func (s *healthService) checkDatabase() *models.ComponentHealth {
	start := time.Now()
	component := &models.ComponentHealth{
		Name:    "database",
		Details: make(map[string]string),
	}
	
	if s.db == nil {
		component.Status = models.HealthStatusUnhealthy
		component.Error = "Database connection not initialized"
		return component
	}
	
	// Create a context with timeout for the health check
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Check database health
	if err := s.db.Health(); err != nil {
		component.Status = models.HealthStatusUnhealthy
		component.Error = fmt.Sprintf("Database health check failed: %v", err)
		component.ResponseTime = int64(time.Since(start).Nanoseconds() / 1e6)
		return component
	}
	
	// Get database statistics
	sqlDB, err := s.db.DB.DB()
	if err == nil {
		stats := sqlDB.Stats()
		component.Details["open_connections"] = fmt.Sprintf("%d", stats.OpenConnections)
		component.Details["in_use"] = fmt.Sprintf("%d", stats.InUse)
		component.Details["idle"] = fmt.Sprintf("%d", stats.Idle)
		component.Details["max_open"] = fmt.Sprintf("%d", stats.MaxOpenConnections)
		
		// Check if we're approaching connection limits
		if stats.OpenConnections > int(float64(stats.MaxOpenConnections)*0.8) {
			component.Status = models.HealthStatusDegraded
			component.Message = "High connection usage detected"
		} else {
			component.Status = models.HealthStatusHealthy
			component.Message = "Database connection healthy"
		}
	} else {
		component.Status = models.HealthStatusDegraded
		component.Message = "Could not get connection stats"
		component.Error = err.Error()
	}
	
	component.ResponseTime = int64(time.Since(start).Nanoseconds() / 1e6)
	return component
}

func (s *healthService) checkMemory() *models.ComponentHealth {
	component := &models.ComponentHealth{
		Name:    "memory",
		Details: make(map[string]string),
	}
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Convert bytes to MB for readability
	allocMB := m.Alloc / 1024 / 1024
	sysMB := m.Sys / 1024 / 1024
	
	component.Details["alloc_mb"] = fmt.Sprintf("%d", allocMB)
	component.Details["sys_mb"] = fmt.Sprintf("%d", sysMB)
	component.Details["num_gc"] = fmt.Sprintf("%d", m.NumGC)
	
	// Simple memory usage assessment
	if allocMB > 512 { // More than 512MB allocated
		component.Status = models.HealthStatusDegraded
		component.Message = "High memory usage detected"
	} else if allocMB > 1024 { // More than 1GB allocated
		component.Status = models.HealthStatusUnhealthy
		component.Message = "Critical memory usage"
	} else {
		component.Status = models.HealthStatusHealthy
		component.Message = "Memory usage normal"
	}
	
	return component
}

func (s *healthService) checkGoroutines() *models.ComponentHealth {
	component := &models.ComponentHealth{
		Name:    "goroutines",
		Details: make(map[string]string),
	}
	
	numGoroutines := runtime.NumGoroutine()
	component.Details["count"] = fmt.Sprintf("%d", numGoroutines)
	
	// Assess goroutine count
	if numGoroutines > 1000 {
		component.Status = models.HealthStatusDegraded
		component.Message = "High goroutine count detected"
	} else if numGoroutines > 5000 {
		component.Status = models.HealthStatusUnhealthy
		component.Message = "Critical goroutine count"
	} else {
		component.Status = models.HealthStatusHealthy
		component.Message = "Goroutine count normal"
	}
	
	return component
}

func (s *healthService) calculateSummary(components []models.ComponentHealth) models.HealthSummary {
	summary := models.HealthSummary{
		TotalComponents: len(components),
	}
	
	for _, comp := range components {
		switch comp.Status {
		case models.HealthStatusHealthy:
			summary.HealthyCount++
		case models.HealthStatusDegraded:
			summary.DegradedCount++
		case models.HealthStatusUnhealthy:
			summary.UnhealthyCount++
		}
	}
	
	return summary
}

func (s *healthService) determineOverallStatus(components []models.ComponentHealth) models.HealthStatus {
	hasUnhealthy := false
	hasDegraded := false
	
	for _, comp := range components {
		switch comp.Status {
		case models.HealthStatusUnhealthy:
			hasUnhealthy = true
		case models.HealthStatusDegraded:
			hasDegraded = true
		}
	}
	
	if hasUnhealthy {
		return models.HealthStatusUnhealthy
	} else if hasDegraded {
		return models.HealthStatusDegraded
	}
	
	return models.HealthStatusHealthy
}