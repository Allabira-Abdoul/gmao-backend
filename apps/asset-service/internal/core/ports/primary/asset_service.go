package primary

import (
	"context"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
)

// AssetService defines primary application business operations.
type AssetService interface {
	CreateAsset(ctx context.Context, req domain.CreateAssetRequest) (*domain.AssetResponse, error)
	UpdateAsset(ctx context.Context, id uuid.UUID, req domain.UpdateAssetRequest) (*domain.AssetResponse, error)
	DeleteAsset(ctx context.Context, id uuid.UUID) error
	GetAsset(ctx context.Context, id uuid.UUID) (*domain.AssetResponse, error)
	GetAssetByCode(ctx context.Context, code string) (*domain.AssetResponse, error)
	GetAllAssets(ctx context.Context) ([]domain.AssetResponse, error)
}
