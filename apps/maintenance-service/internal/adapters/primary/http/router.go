package http

import (
	"backend-gmao/apps/maintenance-service/internal/application"
	"backend-gmao/pkg/auth"
	authdomain "backend-gmao/pkg/auth/domain"
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
			ordres.GET("", middleware.RequirePrivilege(authdomain.PrivilegeWorkOrderView), otHandler.ListOrdresTravail)
			ordres.GET("/:id", middleware.RequirePrivilege(authdomain.PrivilegeWorkOrderView), otHandler.GetOrdreTravail)
			ordres.POST("", middleware.RequirePrivilege(authdomain.PrivilegeWorkOrderCreate), otHandler.CreateOrdreTravail)
			ordres.PUT("/:id", middleware.RequirePrivilege(authdomain.PrivilegeWorkOrderUpdate), otHandler.UpdateOrdreTravail)
			ordres.DELETE("/:id", middleware.RequirePrivilege(authdomain.PrivilegeWorkOrderDelete), otHandler.DeleteOrdreTravail)
		}

		// Intervention CRUD
		interventions := authenticated.Group("/interventions")
		{
			interventions.GET("", middleware.RequirePrivilege(authdomain.PrivilegeMaintenanceView), interventionHandler.ListInterventions)
			interventions.GET("/:id", middleware.RequirePrivilege(authdomain.PrivilegeMaintenanceView), interventionHandler.GetIntervention)
			interventions.POST("", middleware.RequirePrivilege(authdomain.PrivilegeMaintenancePlanCreate), interventionHandler.CreateIntervention)
			interventions.PUT("/:id", middleware.RequirePrivilege(authdomain.PrivilegeMaintenancePlanUpdate), interventionHandler.UpdateIntervention)
			interventions.DELETE("/:id", middleware.RequirePrivilege(authdomain.PrivilegeMaintenancePlanDelete), interventionHandler.DeleteIntervention)
		}
	}
}
