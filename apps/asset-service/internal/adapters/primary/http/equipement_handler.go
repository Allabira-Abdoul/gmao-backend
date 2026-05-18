package http

import (
	"errors"
	"net/http"

	"backend-gmao/apps/asset-service/internal/application"
	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EquipementHandler handles HTTP requests for equipement operations.
type EquipementHandler struct {
	service *application.EquipementService
}

// NewEquipementHandler creates a new EquipementHandler.
func NewEquipementHandler(service *application.EquipementService) *EquipementHandler {
	return &EquipementHandler{service: service}
}

// ListEquipements handles GET /equipements
func (h *EquipementHandler) ListEquipements(c *gin.Context) {
	pagination := response.GetPagination(c, 1, 20)

	equipements, total, err := h.service.ListEquipements(c.Request.Context(), pagination.Limit, pagination.Offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list equipements")
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, equipements, response.NewMeta(pagination.Page, pagination.PerPage, total))
}

// GetEquipement handles GET /equipements/:id
func (h *EquipementHandler) GetEquipement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid equipement ID format")
		return
	}

	equipement, err := h.service.GetEquipementByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, application.ErrEquipementNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Equipement not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get equipement")
		return
	}

	response.Success(c, http.StatusOK, equipement)
}

// CreateEquipement handles POST /equipements
func (h *EquipementHandler) CreateEquipement(c *gin.Context) {
	var req domain.CreateEquipementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	equipement, err := h.service.CreateEquipement(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, application.ErrCodeExists) {
			response.Error(c, http.StatusConflict, "CODE_EXISTS", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create equipement")
		return
	}

	response.Success(c, http.StatusCreated, equipement)
}

// UpdateEquipement handles PUT /equipements/:id
func (h *EquipementHandler) UpdateEquipement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid equipement ID format")
		return
	}

	var req domain.UpdateEquipementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	equipement, err := h.service.UpdateEquipement(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, application.ErrEquipementNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Equipement not found")
			return
		}
		if errors.Is(err, application.ErrCodeExists) {
			response.Error(c, http.StatusConflict, "CODE_EXISTS", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update equipement")
		return
	}

	response.Success(c, http.StatusOK, equipement)
}

// DeleteEquipement handles DELETE /equipements/:id
func (h *EquipementHandler) DeleteEquipement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid equipement ID format")
		return
	}

	if err := h.service.DeleteEquipement(c.Request.Context(), id); err != nil {
		if errors.Is(err, application.ErrEquipementNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Equipement not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete equipement")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Equipement deleted successfully"})
}
