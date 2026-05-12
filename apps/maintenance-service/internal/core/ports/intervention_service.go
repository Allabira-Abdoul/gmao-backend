package ports

import (
	"context"

	"github.com/google/uuid"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
)

// InterventionService defines the use cases for intervention operations.
type InterventionService interface {
	CreateIntervention(ctx context.Context, req *domain.CreateInterventionRequest) (*domain.InterventionResponse, error)
	GetInterventionByID(ctx context.Context, id uuid.UUID) (*domain.InterventionResponse, error)
	ListByOrdreTravail(ctx context.Context, ordreTravailID uuid.UUID, limit, offset int) ([]domain.InterventionResponse, int64, error)
	ListByTechnicien(ctx context.Context, technicienID uuid.UUID, limit, offset int) ([]domain.InterventionResponse, int64, error)
	UpdateIntervention(ctx context.Context, id uuid.UUID, req *domain.UpdateInterventionRequest) (*domain.InterventionResponse, error)
	DeleteIntervention(ctx context.Context, id uuid.UUID) error
}
