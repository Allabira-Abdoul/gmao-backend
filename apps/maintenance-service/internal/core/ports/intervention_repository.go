package ports

import (
	"context"

	"github.com/google/uuid"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
)

// InterventionRepository defines the interface for data access operations related to Intervention.
type InterventionRepository interface {
	Create(ctx context.Context, intervention *domain.Intervention) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Intervention, error)
	ListByOrdreTravail(ctx context.Context, ordreTravailID uuid.UUID, limit, offset int) ([]domain.Intervention, int64, error)
	ListByTechnicien(ctx context.Context, technicienID uuid.UUID, limit, offset int) ([]domain.Intervention, int64, error)
	Update(ctx context.Context, intervention *domain.Intervention) error
	Delete(ctx context.Context, id uuid.UUID) error
}
