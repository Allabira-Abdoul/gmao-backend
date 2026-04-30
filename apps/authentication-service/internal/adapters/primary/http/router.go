package http

import (
	"backend-gmao/apps/authentication-service/internal/application"
	"github.com/gin-gonic/gin"
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
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
	}
}
