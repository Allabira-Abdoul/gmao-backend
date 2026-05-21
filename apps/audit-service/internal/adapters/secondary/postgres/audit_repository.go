package postgres

import (
	"context"

	"backend-gmao/apps/audit-service/internal/core/domain"
	"gorm.io/gorm"
)

type auditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a GORM audit log repository.
func NewAuditRepository(db *gorm.DB) *auditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditRepository) FindAll(ctx context.Context) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	// Order by PerformedAt desc so that the latest events appear first
	if err := r.db.WithContext(ctx).Order("performed_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
