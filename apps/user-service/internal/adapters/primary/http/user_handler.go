package http

import (
	"errors"
	"fmt"
	"net/http"

	"backend-gmao/apps/user-service/internal/application"
	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/pkg/middleware"
	"backend-gmao/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	service *application.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	page := 1
	perPage := 20

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if pp := c.Query("per_page"); pp != "" {
		fmt.Sscanf(pp, "%d", &perPage)
	}

	users, total, err := h.service.ListUsers(c.Request.Context(), page, perPage)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list users")
		return
	}

	totalPages := total / int64(perPage)
	if total%int64(perPage) != 0 {
		totalPages++
	}

	response.SuccessWithMeta(c, http.StatusOK, users, &response.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get user")
		return
	}

	response.Success(c, http.StatusOK, user)
}

// GetCurrentUser handles GET /users/me
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userIDStr, _ := c.Get(middleware.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_TOKEN", "Invalid user ID in token")
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	response.Success(c, http.StatusOK, user)
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, application.ErrEmailExists) {
			response.Error(c, http.StatusConflict, "EMAIL_EXISTS", err.Error())
			return
		}
		if errors.Is(err, application.ErrRoleNotFound) {
			response.Error(c, http.StatusBadRequest, "ROLE_NOT_FOUND", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create user")
		return
	}

	response.Success(c, http.StatusCreated, user)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		if errors.Is(err, application.ErrEmailExists) {
			response.Error(c, http.StatusConflict, "EMAIL_EXISTS", err.Error())
			return
		}
		if errors.Is(err, application.ErrRoleNotFound) {
			response.Error(c, http.StatusBadRequest, "ROLE_NOT_FOUND", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update user")
		return
	}

	response.Success(c, http.StatusOK, user)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		if errors.Is(err, application.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete user")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}
