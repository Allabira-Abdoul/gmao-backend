package http

import (
	"backend-gmao/apps/authentication-service/internal/application"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

// RegisterRoutes sets up all HTTP routes for the authentication service.
func RegisterRoutes(
	router *gin.Engine,
	authService *application.AuthService,
	healthHandler *HealthHandler,
) {
	authHandler := NewAuthHandler(authService)

	// Health check
	router.GET("/health", healthHandler.HealthCheck)

	// Auth group
	auth := router.Group("/auth")
	{
		auth.POST("/login", middleware.RateLimit(5, time.Minute), authHandler.Login)
		auth.POST("/refresh", middleware.RateLimit(20, time.Minute), authHandler.Refresh)
		auth.POST("/logout", middleware.RateLimit(20, time.Minute), authHandler.Logout)
	}
}
