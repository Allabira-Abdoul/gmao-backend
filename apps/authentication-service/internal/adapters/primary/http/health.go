package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	dbStatus := "UP"

	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "DOWN"
		}
	} else {
		dbStatus = "NOT_CONFIGURED"
	}

	status := http.StatusOK
	overallStatus := "UP"
	if dbStatus != "UP" {
		status = http.StatusServiceUnavailable
		overallStatus = "DEGRADED"
	}

	c.JSON(status, gin.H{
		"status":  overallStatus,
		"service": "authentication-service",
		"checks": gin.H{
			"database": dbStatus,
		},
	})
}
