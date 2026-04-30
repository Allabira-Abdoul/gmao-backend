package http

import (
	"backend-gmao/apps/user-service/internal/application"
	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the user service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	userService *application.UserService,
	roleService *application.RoleService,
) {
	userHandler := NewUserHandler(userService)
	roleHandler := NewRoleHandler(roleService)
	internalHandler := NewInternalHandler(userService)

	// --- Internal endpoints (service-to-service only) ---
	internal := router.Group("/internal")
	internal.Use(middleware.RequireInternalService())
	{
		internal.GET("/by-email", internalHandler.GetUserByEmail)
		internal.GET("/by-id", internalHandler.GetUserByID)
	}

	// --- Public endpoints ---
	// System privileges reference (useful for admin UIs)
	router.GET("/privileges", func(c *gin.Context) {
		roleHandler.ListPrivileges(c)
	})

	// --- Authenticated endpoints ---
	authenticated := router.Group("/")
	authenticated.Use(middleware.RequireAuth(jwtManager))
	{
		// Current user profile (any authenticated user)
		authenticated.GET("/users/me", userHandler.GetCurrentUser)

		// User CRUD (privilege-protected)
		users := authenticated.Group("/users")
		{
			users.GET("", middleware.RequirePrivilege(domain.PrivilegeUserView), userHandler.ListUsers)
			users.GET("/:id", middleware.RequirePrivilege(domain.PrivilegeUserView), userHandler.GetUser)
			users.POST("", middleware.RequirePrivilege(domain.PrivilegeUserCreate), userHandler.CreateUser)
			users.PUT("/:id", middleware.RequirePrivilege(domain.PrivilegeUserUpdate), userHandler.UpdateUser)
			users.DELETE("/:id", middleware.RequirePrivilege(domain.PrivilegeUserDelete), userHandler.DeleteUser)
		}

		// Role CRUD (privilege-protected)
		roles := authenticated.Group("/roles")
		{
			roles.GET("", middleware.RequirePrivilege(domain.PrivilegeRoleView), roleHandler.ListRoles)
			roles.GET("/:id", middleware.RequirePrivilege(domain.PrivilegeRoleView), roleHandler.GetRole)
			roles.POST("", middleware.RequirePrivilege(domain.PrivilegeRoleCreate), roleHandler.CreateRole)
			roles.PUT("/:id", middleware.RequirePrivilege(domain.PrivilegeRoleUpdate), roleHandler.UpdateRole)
			roles.DELETE("/:id", middleware.RequirePrivilege(domain.PrivilegeRoleDelete), roleHandler.DeleteRole)
			roles.PUT("/:id/privileges", middleware.RequirePrivilege(domain.PrivilegeRoleUpdate), roleHandler.SetRolePrivileges)
		}
	}
}
