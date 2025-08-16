package models

import (
	"time"
)

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

type ComponentHealth struct {
	Name         string            `json:"name" example:"database"`
	Status       HealthStatus      `json:"status" example:"healthy"`
	Message      string            `json:"message,omitempty" example:"Connection successful"`
	ResponseTime int64             `json:"response_time_ms" example:"15"`
	Details      map[string]string `json:"details,omitempty"`
	Error        string            `json:"error,omitempty" example:"connection timeout"`
}

type HealthResponse struct {
	Status      HealthStatus      `json:"status" example:"healthy"`
	Timestamp   time.Time         `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Version     string            `json:"version" example:"1.0.0"`
	Uptime      int64             `json:"uptime_seconds" example:"3600"`
	Components  []ComponentHealth `json:"components"`
	Summary     HealthSummary     `json:"summary"`
}

type HealthSummary struct {
	TotalComponents int `json:"total_components" example:"3"`
	HealthyCount    int `json:"healthy" example:"2"`
	DegradedCount   int `json:"degraded" example:"1"`
	UnhealthyCount  int `json:"unhealthy" example:"0"`
}

type LivenessResponse struct {
	Status    HealthStatus `json:"status" example:"healthy"`
	Timestamp time.Time    `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Message   string       `json:"message" example:"Service is alive"`
}

type ReadinessResponse struct {
	Status     HealthStatus      `json:"status" example:"healthy"`
	Timestamp  time.Time         `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Message    string            `json:"message" example:"Service is ready"`
	Components []ComponentHealth `json:"components"`
}