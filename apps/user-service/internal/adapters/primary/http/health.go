package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck provides a simple health check response
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"service": "user-service",
	})
}
