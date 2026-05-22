package http

import (
	"backend-gmao/apps/asset-service/internal/application/service"
	"backend-gmao/pkg/auth"
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
			assets.POST("", middleware.RequirePrivilege("ASSET_CREATE"), assetHandler.CreateAsset)
			assets.GET("", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.ListAssets)
			assets.GET("/:id", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetAsset)
			assets.GET("/code/:code", middleware.RequirePrivilege("ASSET_VIEW"), assetHandler.GetAssetByCode)
			assets.PUT("/:id", middleware.RequirePrivilege("ASSET_UPDATE"), assetHandler.UpdateAsset)
			assets.DELETE("/:id", middleware.RequirePrivilege("ASSET_DELETE"), assetHandler.DeleteAsset)
		}
	}
}
