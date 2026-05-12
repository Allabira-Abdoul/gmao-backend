package http

import (
	"backend-gmao/apps/maintenance-service/internal/application"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the maintenance service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	ordreTravailService *application.OrdreTravailService,
	interventionService *application.InterventionService,
) {
	otHandler := NewOrdreTravailHandler(ordreTravailService)
	interventionHandler := NewInterventionHandler(interventionService)

	// --- Authenticated endpoints ---
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		// Work Order CRUD
		ordres := authenticated.Group("/ordres-travail")
		{
			ordres.GET("", otHandler.ListOrdresTravail)
			ordres.GET("/:id", otHandler.GetOrdreTravail)
			ordres.POST("", otHandler.CreateOrdreTravail)
			ordres.PUT("/:id", otHandler.UpdateOrdreTravail)
			ordres.DELETE("/:id", otHandler.DeleteOrdreTravail)
		}

		// Intervention CRUD
		interventions := authenticated.Group("/interventions")
		{
			interventions.GET("", interventionHandler.ListInterventions)
			interventions.GET("/:id", interventionHandler.GetIntervention)
			interventions.POST("", interventionHandler.CreateIntervention)
			interventions.PUT("/:id", interventionHandler.UpdateIntervention)
			interventions.DELETE("/:id", interventionHandler.DeleteIntervention)
		}
	}
}
