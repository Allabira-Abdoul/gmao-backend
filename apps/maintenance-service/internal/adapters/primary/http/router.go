package http

import (
	"backend-gmao/apps/maintenance-service/internal/application/service"
	"backend-gmao/pkg/auth"
	authdomain "backend-gmao/pkg/auth/domain"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the maintenance service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	maintenanceService *service.MaintenanceService,
) {
	maintenanceHandler := NewMaintenanceHandler(maintenanceService)

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		workorders := authenticated.Group("/work-orders")
		{
<<<<<<< HEAD
			workorders.POST("", maintenanceHandler.CreateWorkOrder)
			workorders.GET("", maintenanceHandler.ListWorkOrders)
			workorders.GET("/:id", maintenanceHandler.GetWorkOrder)
			workorders.PUT("/:id", maintenanceHandler.UpdateWorkOrder)
			workorders.DELETE("/:id", maintenanceHandler.DeleteWorkOrder)

			// Interventions under work order
			workorders.POST("/:id/interventions", maintenanceHandler.CreateIntervention)
			workorders.GET("/:id/interventions", maintenanceHandler.GetInterventions)
=======
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
>>>>>>> 0240860163106025f1377fcc80aa118b977be644
		}
	}
}
