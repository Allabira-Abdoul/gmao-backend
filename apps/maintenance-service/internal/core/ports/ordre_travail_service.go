package ports

import (
	"context"

	"github.com/google/uuid"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
)

// OrdreTravailService defines the use cases for work order operations.
type OrdreTravailService interface {
	CreateOrdreTravail(ctx context.Context, req *domain.CreateOrdreTravailRequest, createdBy uuid.UUID) (*domain.OrdreTravailResponse, error)
	GetOrdreTravailByID(ctx context.Context, id uuid.UUID) (*domain.OrdreTravailResponse, error)
	ListOrdresTravail(ctx context.Context, limit, offset int) ([]domain.OrdreTravailResponse, int64, error)
	ListByEquipement(ctx context.Context, equipementID uuid.UUID, limit, offset int) ([]domain.OrdreTravailResponse, int64, error)
	ListByStatut(ctx context.Context, statut domain.OrdreTravailStatut, limit, offset int) ([]domain.OrdreTravailResponse, int64, error)
	ListByAssigne(ctx context.Context, utilisateurID uuid.UUID, limit, offset int) ([]domain.OrdreTravailResponse, int64, error)
	UpdateOrdreTravail(ctx context.Context, id uuid.UUID, req *domain.UpdateOrdreTravailRequest) (*domain.OrdreTravailResponse, error)
	DeleteOrdreTravail(ctx context.Context, id uuid.UUID) error
}
