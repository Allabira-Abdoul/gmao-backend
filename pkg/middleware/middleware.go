package middleware

import (
	"net/http"
	"strings"

	"backend-gmao/pkg/auth"
	"github.com/gin-gonic/gin"
)

const (
	// Context keys for storing authenticated user information.
	ContextKeyUserID     = "auth_user_id"
	ContextKeyEmail      = "auth_email"
	ContextKeyRole       = "auth_role"
	ContextKeyPrivileges = "auth_privileges"
)

// RequireAuth returns a Gin middleware that validates the JWT access token
// from the Authorization header and injects claims into the Gin context.
func RequireAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authorization header is required",
				},
			})
			return
		}

		// Expect "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Authorization header must be in the format: Bearer <token>",
				},
			})
			return
		}

		claims, err := jwtManager.ValidateAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired access token",
				},
			})
			return
		}

		// Inject user info into Gin context for downstream handlers
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyRole, claims.Role)
		c.Set(ContextKeyPrivileges, claims.Privileges)

		c.Next()
	}
}

// RequirePrivilege returns a Gin middleware that checks if the authenticated user
// has the specified privilege. Must be used after RequireAuth.
func RequirePrivilege(privilege string) gin.HandlerFunc {
	return func(c *gin.Context) {
		privileges, exists := c.Get(ContextKeyPrivileges)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Access denied: no privileges found",
				},
			})
			return
		}

		userPrivileges, ok := privileges.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to parse user privileges",
				},
			})
			return
		}

		// SYSTEM_ADMIN has access to everything
		for _, p := range userPrivileges {
			if p == "SYSTEM_ADMIN" || p == privilege {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INSUFFICIENT_PRIVILEGES",
				"message": "You do not have the required privilege: " + privilege,
			},
		})
	}
}

// RequireInternalService returns a middleware that ensures the request comes from
// an internal service (via the API gateway or direct inter-service call).
func RequireInternalService() gin.HandlerFunc {
	return func(c *gin.Context) {
		gatewayHeader := c.GetHeader("X-Gateway-Service")
		internalHeader := c.GetHeader("X-Internal-Service")

		if gatewayHeader == "" && internalHeader == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "This endpoint is only accessible internally",
				},
			})
			return
		}

		c.Next()
	}
}
