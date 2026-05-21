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
			assets.POST("", assetHandler.CreateAsset)
			assets.GET("", assetHandler.ListAssets)
			assets.GET("/:id", assetHandler.GetAsset)
			assets.GET("/code/:code", assetHandler.GetAssetByCode)
			assets.PUT("/:id", assetHandler.UpdateAsset)
			assets.DELETE("/:id", assetHandler.DeleteAsset)
		}
	}
}
