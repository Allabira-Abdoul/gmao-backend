package http

import (
	"backend-gmao/apps/asset-service/internal/application"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the asset service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	equipementService *application.EquipementService,
) {
	equipementHandler := NewEquipementHandler(equipementService)

	// --- Authenticated endpoints ---
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		// Equipement CRUD
		equipements := authenticated.Group("/equipements")
		{
			equipements.GET("", equipementHandler.ListEquipements)
			equipements.GET("/:id", equipementHandler.GetEquipement)
			equipements.POST("", equipementHandler.CreateEquipement)
			equipements.PUT("/:id", equipementHandler.UpdateEquipement)
			equipements.DELETE("/:id", equipementHandler.DeleteEquipement)
		}
	}
}
