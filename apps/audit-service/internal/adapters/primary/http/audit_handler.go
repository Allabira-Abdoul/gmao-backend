package http

import (
	"net/http"

	"backend-gmao/apps/audit-service/internal/core/domain"
	"backend-gmao/apps/audit-service/internal/core/ports/primary"
	"github.com/gin-gonic/gin"
)

// AuditHandler manages HTTP routes for audit logs.
type AuditHandler struct {
	auditService primary.AuditService
}

// NewAuditHandler creates a new AuditHandler.
func NewAuditHandler(auditService primary.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// WriteLog records a new audit event log. Only called internally.
func (h *AuditHandler) WriteLog(c *gin.Context) {
	var req domain.CreateAuditLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.auditService.WriteLog(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ListLogs lists all recorded audit events. Only allowed for auditors.
func (h *AuditHandler) ListLogs(c *gin.Context) {
	resp, err := h.auditService.GetAllLogs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
