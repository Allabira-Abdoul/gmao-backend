package postgres

import (
	"context"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type assetRepository struct {
	db *gorm.DB
}

// NewAssetRepository creates a GORM asset repository.
func NewAssetRepository(db *gorm.DB) *assetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) Create(ctx context.Context, asset *domain.Asset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *assetRepository) Update(ctx context.Context, asset *domain.Asset) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *assetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Asset{}, "id = ?", id).Error
}

func (r *assetRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.db.WithContext(ctx).First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) FindByCode(ctx context.Context, code string) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.db.WithContext(ctx).First(&asset, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) FindAll(ctx context.Context) ([]domain.Asset, error) {
	var assets []domain.Asset
	if err := r.db.WithContext(ctx).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}
