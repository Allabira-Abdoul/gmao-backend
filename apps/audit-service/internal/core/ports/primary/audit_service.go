package primary

import (
	"context"

	"backend-gmao/apps/audit-service/internal/core/domain"
)

// AuditService defines primary application business operations.
type AuditService interface {
	WriteLog(ctx context.Context, req domain.CreateAuditLogRequest) (*domain.AuditLogResponse, error)
	GetAllLogs(ctx context.Context) ([]domain.AuditLogResponse, error)
}
