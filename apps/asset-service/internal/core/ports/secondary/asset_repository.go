package secondary

import (
	"context"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
)

// AssetRepository defines secondary adapter database actions.
type AssetRepository interface {
	Create(ctx context.Context, asset *domain.Asset) error
	Update(ctx context.Context, asset *domain.Asset) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
	FindByCode(ctx context.Context, code string) (*domain.Asset, error)
	FindAll(ctx context.Context) ([]domain.Asset, error)
}
