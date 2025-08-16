package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"postService/internal/models"
	"postService/internal/service"
)

type HealthHandler struct {
	healthService service.HealthService
}

func NewHealthHandler(healthService service.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// GetHealth godoc
// @Summary Get comprehensive health status
// @Description Get detailed health information including all component statuses
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse "Service is healthy"
// @Success 503 {object} models.HealthResponse "Service is degraded or unhealthy"
// @Router /health [get]
func (h *HealthHandler) GetHealth(c *gin.Context) {
	health := h.healthService.GetHealth()
	
	// Set appropriate HTTP status code based on health
	var statusCode int
	switch health.Status {
	case models.HealthStatusHealthy:
		statusCode = http.StatusOK
	case models.HealthStatusDegraded:
		statusCode = http.StatusServiceUnavailable
	case models.HealthStatusUnhealthy:
		statusCode = http.StatusServiceUnavailable
	default:
		statusCode = http.StatusServiceUnavailable
	}
	
	c.JSON(statusCode, health)
}

// GetLiveness godoc
// @Summary Liveness probe endpoint
// @Description Kubernetes liveness probe - indicates if the service is alive
// @Tags health
// @Produce json
// @Success 200 {object} models.LivenessResponse "Service is alive"
// @Router /health/live [get]
func (h *HealthHandler) GetLiveness(c *gin.Context) {
	liveness := h.healthService.GetLiveness()
	c.JSON(http.StatusOK, liveness)
}

// GetReadiness godoc
// @Summary Readiness probe endpoint
// @Description Kubernetes readiness probe - indicates if the service is ready to accept traffic
// @Tags health
// @Produce json
// @Success 200 {object} models.ReadinessResponse "Service is ready"
// @Success 503 {object} models.ReadinessResponse "Service is not ready"
// @Router /health/ready [get]
func (h *HealthHandler) GetReadiness(c *gin.Context) {
	readiness := h.healthService.GetReadiness()
	
	// Set appropriate HTTP status code based on readiness
	var statusCode int
	switch readiness.Status {
	case models.HealthStatusHealthy:
		statusCode = http.StatusOK
	case models.HealthStatusDegraded:
		statusCode = http.StatusServiceUnavailable
	case models.HealthStatusUnhealthy:
		statusCode = http.StatusServiceUnavailable
	default:
		statusCode = http.StatusServiceUnavailable
	}
	
	c.JSON(statusCode, readiness)
}

// GetComponentHealth godoc
// @Summary Get specific component health
// @Description Get health status of a specific component
// @Tags health
// @Produce json
// @Param component path string true "Component name (database, memory, goroutines)"
// @Success 200 {object} models.ComponentHealth "Component is healthy"
// @Success 503 {object} models.ComponentHealth "Component is degraded or unhealthy"
// @Success 404 {object} map[string]string "Component not found"
// @Router /health/component/{component} [get]
func (h *HealthHandler) GetComponentHealth(c *gin.Context) {
	componentName := c.Param("component")
	
	if componentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Component name is required",
		})
		return
	}
	
	componentHealth := h.healthService.CheckComponent(componentName)
	
	// Set appropriate HTTP status code based on component health
	var statusCode int
	switch componentHealth.Status {
	case models.HealthStatusHealthy:
		statusCode = http.StatusOK
	case models.HealthStatusDegraded:
		statusCode = http.StatusServiceUnavailable
	case models.HealthStatusUnhealthy:
		statusCode = http.StatusServiceUnavailable
	default:
		statusCode = http.StatusNotFound
	}
	
	c.JSON(statusCode, componentHealth)
}

// GetHealthSimple godoc
// @Summary Simple health check
// @Description Simple health check that returns OK if service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "Service is OK"
// @Router /health/ping [get]
func (h *HealthHandler) GetHealthSimple(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}