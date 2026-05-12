package ports

import (
	"context"

	"github.com/google/uuid"

	"backend-gmao/apps/asset-service/internal/core/domain"
)

// EquipementRepository defines the interface for data access operations related to Equipement.
type EquipementRepository interface {
	Create(ctx context.Context, equipement *domain.Equipement) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Equipement, error)
	GetByCode(ctx context.Context, code string) (*domain.Equipement, error)
	List(ctx context.Context, limit, offset int) ([]domain.Equipement, int64, error)
	Update(ctx context.Context, equipement *domain.Equipement) error
	Delete(ctx context.Context, id uuid.UUID) error
}
