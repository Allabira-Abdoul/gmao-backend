package secondary

import (
	"context"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"github.com/google/uuid"
)

// MaintenanceRepository defines secondary adapter database actions.
type MaintenanceRepository interface {
	CreateWorkOrder(ctx context.Context, wo *domain.OrdreTravail) error
	UpdateWorkOrder(ctx context.Context, wo *domain.OrdreTravail) error
	DeleteWorkOrder(ctx context.Context, id uuid.UUID) error
	FindWorkOrderByID(ctx context.Context, id uuid.UUID) (*domain.OrdreTravail, error)
	FindAllWorkOrders(ctx context.Context) ([]domain.OrdreTravail, error)

	CreateIntervention(ctx context.Context, intervention *domain.Intervention) error
	FindInterventionsByWorkOrderID(ctx context.Context, workOrderID uuid.UUID) ([]domain.Intervention, error)
}
