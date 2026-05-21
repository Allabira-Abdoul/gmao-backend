package http

import (
	"backend-gmao/apps/analytics-service/internal/application/service"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the analytics service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	analyticsService *service.AnalyticsService,
) {
	metricHandler := NewMetricHandler(analyticsService)

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		metrics := authenticated.Group("/metrics")
		{
			metrics.POST("", metricHandler.RecordMetric)
			metrics.GET("", metricHandler.ListMetrics)
			metrics.GET("/:id", metricHandler.GetMetric)
			metrics.GET("/category/:category", metricHandler.ListMetricsByCategory)
		}
	}
}
