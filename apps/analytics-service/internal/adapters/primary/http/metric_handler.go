package http

import (
	"net/http"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"backend-gmao/apps/analytics-service/internal/core/ports/primary"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MetricHandler manages HTTP routes for metrics.
type MetricHandler struct {
	analyticsService primary.AnalyticsService
}

// NewMetricHandler creates a new MetricHandler.
func NewMetricHandler(analyticsService primary.AnalyticsService) *MetricHandler {
	return &MetricHandler{analyticsService: analyticsService}
}

// RecordMetric handles metric creation.
func (h *MetricHandler) RecordMetric(c *gin.Context) {
	var req domain.CreateMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.analyticsService.RecordMetric(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetMetric gets a specific metric.
func (h *MetricHandler) GetMetric(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	resp, err := h.analyticsService.GetMetric(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListMetrics lists all metrics.
func (h *MetricHandler) ListMetrics(c *gin.Context) {
	resp, err := h.analyticsService.GetAllMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListMetricsByCategory lists metrics filtering by category name.
func (h *MetricHandler) ListMetricsByCategory(c *gin.Context) {
	category := c.Param("category")
	resp, err := h.analyticsService.GetMetricsByCategory(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
