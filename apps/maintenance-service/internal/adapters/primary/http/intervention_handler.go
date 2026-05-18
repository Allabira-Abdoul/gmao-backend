package http

import (
	"errors"
	"net/http"

	"backend-gmao/apps/maintenance-service/internal/application"
	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InterventionHandler handles HTTP requests for intervention operations.
type InterventionHandler struct {
	service *application.InterventionService
}

// NewInterventionHandler creates a new InterventionHandler.
func NewInterventionHandler(service *application.InterventionService) *InterventionHandler {
	return &InterventionHandler{service: service}
}

// ListInterventions handles GET /interventions
func (h *InterventionHandler) ListInterventions(c *gin.Context) {
	pagination := response.GetPagination(c, 1, 20)

	// Filter by work order
	if otID := c.Query("ordre_travail_id"); otID != "" {
		id, err := uuid.Parse(otID)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid work order ID format")
			return
		}
		interventions, total, err := h.service.ListByOrdreTravail(c.Request.Context(), id, pagination.Limit, pagination.Offset)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list interventions")
			return
		}
		response.SuccessWithMeta(c, http.StatusOK, interventions, response.NewMeta(pagination.Page, pagination.PerPage, total))
		return
	}

	// Filter by technician
	if techID := c.Query("technicien_id"); techID != "" {
		id, err := uuid.Parse(techID)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid technician ID format")
			return
		}
		interventions, total, err := h.service.ListByTechnicien(c.Request.Context(), id, pagination.Limit, pagination.Offset)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list interventions")
			return
		}
		response.SuccessWithMeta(c, http.StatusOK, interventions, response.NewMeta(pagination.Page, pagination.PerPage, total))
		return
	}

	// Without a filter, return an error (interventions must be scoped)
	response.Error(c, http.StatusBadRequest, "MISSING_FILTER", "Please provide 'ordre_travail_id' or 'technicien_id' query parameter")
}

// GetIntervention handles GET /interventions/:id
func (h *InterventionHandler) GetIntervention(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid intervention ID format")
		return
	}

	intervention, err := h.service.GetInterventionByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, application.ErrInterventionNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Intervention not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get intervention")
		return
	}

	response.Success(c, http.StatusOK, intervention)
}

// CreateIntervention handles POST /interventions
func (h *InterventionHandler) CreateIntervention(c *gin.Context) {
	var req domain.CreateInterventionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	intervention, err := h.service.CreateIntervention(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create intervention")
		return
	}

	response.Success(c, http.StatusCreated, intervention)
}

// UpdateIntervention handles PUT /interventions/:id
func (h *InterventionHandler) UpdateIntervention(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid intervention ID format")
		return
	}

	var req domain.UpdateInterventionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	intervention, err := h.service.UpdateIntervention(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, application.ErrInterventionNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Intervention not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update intervention")
		return
	}

	response.Success(c, http.StatusOK, intervention)
}

// DeleteIntervention handles DELETE /interventions/:id
func (h *InterventionHandler) DeleteIntervention(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid intervention ID format")
		return
	}

	if err := h.service.DeleteIntervention(c.Request.Context(), id); err != nil {
		if errors.Is(err, application.ErrInterventionNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Intervention not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete intervention")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Intervention deleted successfully"})
}
