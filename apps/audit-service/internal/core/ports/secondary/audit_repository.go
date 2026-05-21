package secondary

import (
	"context"

	"backend-gmao/apps/audit-service/internal/core/domain"
)

// AuditLogRepository defines secondary adapter database actions.
// Crucially, it excludes modification and deletion interfaces to guarantee append-only immutability.
type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	FindAll(ctx context.Context) ([]domain.AuditLog, error)
}
