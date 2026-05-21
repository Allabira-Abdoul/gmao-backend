package http

import (
	"backend-gmao/apps/asset-service/internal/application/service"
	"backend-gmao/pkg/auth"
	authdomain "backend-gmao/pkg/auth/domain"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the asset service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	assetService *service.AssetService,
) {
	assetHandler := NewAssetHandler(assetService)

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		assets := authenticated.Group("/assets")
		{
<<<<<<< HEAD
			assets.POST("", assetHandler.CreateAsset)
			assets.GET("", assetHandler.ListAssets)
			assets.GET("/:id", assetHandler.GetAsset)
			assets.GET("/code/:code", assetHandler.GetAssetByCode)
			assets.PUT("/:id", assetHandler.UpdateAsset)
			assets.DELETE("/:id", assetHandler.DeleteAsset)
=======
			equipements.GET("", middleware.RequirePrivilege(authdomain.PrivilegeAssetView), equipementHandler.ListEquipements)
			equipements.GET("/:id", middleware.RequirePrivilege(authdomain.PrivilegeAssetView), equipementHandler.GetEquipement)
			equipements.POST("", middleware.RequirePrivilege(authdomain.PrivilegeAssetCreate), equipementHandler.CreateEquipement)
			equipements.PUT("/:id", middleware.RequirePrivilege(authdomain.PrivilegeAssetUpdate), equipementHandler.UpdateEquipement)
			equipements.DELETE("/:id", middleware.RequirePrivilege(authdomain.PrivilegeAssetDelete), equipementHandler.DeleteEquipement)
>>>>>>> 0240860163106025f1377fcc80aa118b977be644
		}
	}
}
