package http

import (
	"errors"
	"net/http"

	"backend-gmao/apps/maintenance-service/internal/application"
	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/pkg/middleware"
	"backend-gmao/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrdreTravailHandler handles HTTP requests for work order operations.
type OrdreTravailHandler struct {
	service *application.OrdreTravailService
}

// NewOrdreTravailHandler creates a new OrdreTravailHandler.
func NewOrdreTravailHandler(service *application.OrdreTravailService) *OrdreTravailHandler {
	return &OrdreTravailHandler{service: service}
}

// ListOrdresTravail handles GET /ordres-travail
func (h *OrdreTravailHandler) ListOrdresTravail(c *gin.Context) {
	pagination := response.GetPagination(c, 1, 20)

	// Check for optional filters
	if equipementID := c.Query("equipement_id"); equipementID != "" {
		id, err := uuid.Parse(equipementID)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid equipement ID format")
			return
		}
		ordres, total, err := h.service.ListByEquipement(c.Request.Context(), id, pagination.Limit, pagination.Offset)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list work orders")
			return
		}
		response.SuccessWithMeta(c, http.StatusOK, ordres, response.NewMeta(pagination.Page, pagination.PerPage, total))
		return
	}

	if statut := c.Query("statut"); statut != "" {
		ordres, total, err := h.service.ListByStatut(c.Request.Context(), domain.OrdreTravailStatut(statut), pagination.Limit, pagination.Offset)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list work orders")
			return
		}
		response.SuccessWithMeta(c, http.StatusOK, ordres, response.NewMeta(pagination.Page, pagination.PerPage, total))
		return
	}

	if assigneID := c.Query("assigne_id"); assigneID != "" {
		id, err := uuid.Parse(assigneID)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
			return
		}
		ordres, total, err := h.service.ListByAssigne(c.Request.Context(), id, pagination.Limit, pagination.Offset)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list work orders")
			return
		}
		response.SuccessWithMeta(c, http.StatusOK, ordres, response.NewMeta(pagination.Page, pagination.PerPage, total))
		return
	}

	ordres, total, err := h.service.ListOrdresTravail(c.Request.Context(), pagination.Limit, pagination.Offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list work orders")
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, ordres, response.NewMeta(pagination.Page, pagination.PerPage, total))
}

// GetOrdreTravail handles GET /ordres-travail/:id
func (h *OrdreTravailHandler) GetOrdreTravail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid work order ID format")
		return
	}

	ordre, err := h.service.GetOrdreTravailByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, application.ErrOrdreTravailNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Work order not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get work order")
		return
	}

	response.Success(c, http.StatusOK, ordre)
}

// CreateOrdreTravail handles POST /ordres-travail
func (h *OrdreTravailHandler) CreateOrdreTravail(c *gin.Context) {
	var req domain.CreateOrdreTravailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get creator ID from JWT context
	userIDStr, _ := c.Get(middleware.ContextKeyUserID)
	createdBy, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_TOKEN", "Invalid user ID in token")
		return
	}

	ordre, err := h.service.CreateOrdreTravail(c.Request.Context(), &req, createdBy)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create work order")
		return
	}

	response.Success(c, http.StatusCreated, ordre)
}

// UpdateOrdreTravail handles PUT /ordres-travail/:id
func (h *OrdreTravailHandler) UpdateOrdreTravail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid work order ID format")
		return
	}

	var req domain.UpdateOrdreTravailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	ordre, err := h.service.UpdateOrdreTravail(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, application.ErrOrdreTravailNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Work order not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update work order")
		return
	}

	response.Success(c, http.StatusOK, ordre)
}

// DeleteOrdreTravail handles DELETE /ordres-travail/:id
func (h *OrdreTravailHandler) DeleteOrdreTravail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid work order ID format")
		return
	}

	if err := h.service.DeleteOrdreTravail(c.Request.Context(), id); err != nil {
		if errors.Is(err, application.ErrOrdreTravailNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Work order not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete work order")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Work order deleted successfully"})
}
