package http

import (
	"backend-gmao/apps/asset-service/internal/application"
	"backend-gmao/pkg/auth"
	authdomain "backend-gmao/pkg/auth/domain"
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
			equipements.GET("", middleware.RequirePrivilege(authdomain.PrivilegeAssetView), equipementHandler.ListEquipements)
			equipements.GET("/:id", middleware.RequirePrivilege(authdomain.PrivilegeAssetView), equipementHandler.GetEquipement)
			equipements.POST("", middleware.RequirePrivilege(authdomain.PrivilegeAssetCreate), equipementHandler.CreateEquipement)
			equipements.PUT("/:id", middleware.RequirePrivilege(authdomain.PrivilegeAssetUpdate), equipementHandler.UpdateEquipement)
			equipements.DELETE("/:id", middleware.RequirePrivilege(authdomain.PrivilegeAssetDelete), equipementHandler.DeleteEquipement)
		}
	}
}
