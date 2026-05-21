package http

import (
<<<<<<< HEAD
	"backend-gmao/apps/user-service/internal/application/service"
	"backend-gmao/apps/user-service/internal/core/domain"
=======
	"backend-gmao/apps/user-service/internal/application"
>>>>>>> 0240860163106025f1377fcc80aa118b977be644
	"backend-gmao/pkg/auth"
	authdomain "backend-gmao/pkg/auth/domain"
	"backend-gmao/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the user service.
func RegisterRoutes(
	router *gin.Engine,
	jwtManager *auth.JWTManager,
	userService *service.UserService,
	roleService *service.RoleService,
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
			users.GET("", middleware.RequirePrivilege(authdomain.PrivilegeUserView), userHandler.ListUsers)
			users.GET("/:id", middleware.RequirePrivilege(authdomain.PrivilegeUserView), userHandler.GetUser)
			users.POST("", middleware.RequirePrivilege(authdomain.PrivilegeUserCreate), userHandler.CreateUser)
			users.PUT("/:id", middleware.RequirePrivilege(authdomain.PrivilegeUserUpdate), userHandler.UpdateUser)
			users.DELETE("/:id", middleware.RequirePrivilege(authdomain.PrivilegeUserDelete), userHandler.DeleteUser)
		}

		// Role CRUD (privilege-protected)
		roles := authenticated.Group("/roles")
		{
			roles.GET("", middleware.RequirePrivilege(authdomain.PrivilegeRoleView), roleHandler.ListRoles)
			roles.GET("/:id", middleware.RequirePrivilege(authdomain.PrivilegeRoleView), roleHandler.GetRole)
			roles.POST("", middleware.RequirePrivilege(authdomain.PrivilegeRoleCreate), roleHandler.CreateRole)
			roles.PUT("/:id", middleware.RequirePrivilege(authdomain.PrivilegeRoleUpdate), roleHandler.UpdateRole)
			roles.DELETE("/:id", middleware.RequirePrivilege(authdomain.PrivilegeRoleDelete), roleHandler.DeleteRole)
			roles.PUT("/:id/privileges", middleware.RequirePrivilege(authdomain.PrivilegeRoleUpdate), roleHandler.SetRolePrivileges)
		}
	}
}
