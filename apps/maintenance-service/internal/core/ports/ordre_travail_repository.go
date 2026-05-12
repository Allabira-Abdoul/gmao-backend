package ports

import (
	"context"

	"github.com/google/uuid"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
)

// OrdreTravailRepository defines the interface for data access operations related to OrdreTravail.
type OrdreTravailRepository interface {
	Create(ctx context.Context, ordre *domain.OrdreTravail) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.OrdreTravail, error)
	GetByReference(ctx context.Context, reference string) (*domain.OrdreTravail, error)
	List(ctx context.Context, limit, offset int) ([]domain.OrdreTravail, int64, error)
	ListByEquipement(ctx context.Context, equipementID uuid.UUID, limit, offset int) ([]domain.OrdreTravail, int64, error)
	ListByStatut(ctx context.Context, statut domain.OrdreTravailStatut, limit, offset int) ([]domain.OrdreTravail, int64, error)
	ListByAssigne(ctx context.Context, utilisateurID uuid.UUID, limit, offset int) ([]domain.OrdreTravail, int64, error)
	Update(ctx context.Context, ordre *domain.OrdreTravail) error
	Delete(ctx context.Context, id uuid.UUID) error
}
