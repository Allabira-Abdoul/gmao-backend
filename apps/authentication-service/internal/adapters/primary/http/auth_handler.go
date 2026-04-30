package http

import (
	"errors"
	"net/http"

	"backend-gmao/apps/authentication-service/internal/application"
	"backend-gmao/apps/authentication-service/internal/core/domain"
	"backend-gmao/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *application.AuthService
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, application.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
			return
		}
		if errors.Is(err, application.ErrAccountDisabled) {
			response.Error(c, http.StatusForbidden, "ACCOUNT_DISABLED", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to login")
		return
	}

	response.Success(c, http.StatusOK, tokens)
}

// Refresh handles POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	tokens, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, application.ErrInvalidToken) {
			response.Error(c, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to refresh token")
		return
	}

	response.Success(c, http.StatusOK, tokens)
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req domain.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.service.Logout(c.Request.Context(), req); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to logout")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}
