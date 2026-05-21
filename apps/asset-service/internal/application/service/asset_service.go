package service

import (
	"context"
	"errors"
	"time"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/apps/asset-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

var (
	ErrAssetNotFound = errors.New("asset not found")
	ErrCodeExists    = errors.New("an asset with this code already exists")
)

// AssetService implements primary.AssetService.
type AssetService struct {
	assetRepo secondary.AssetRepository
}

// NewAssetService initializes a new AssetService instance.
func NewAssetService(assetRepo secondary.AssetRepository) *AssetService {
	return &AssetService{assetRepo: assetRepo}
}

func (s *AssetService) CreateAsset(ctx context.Context, req domain.CreateAssetRequest) (*domain.AssetResponse, error) {
	existing, _ := s.assetRepo.FindByCode(ctx, req.Code)
	if existing != nil {
		return nil, ErrCodeExists
	}

	asset := &domain.Asset{
		ID:            uuid.New(),
		Name:          req.Name,
		Code:          req.Code,
		Status:        "OPERATIONAL",
		Category:      req.Category,
		Location:      req.Location,
		PurchaseDate:  req.PurchaseDate,
		PurchaseValue: req.PurchaseValue,
	}

	if err := s.assetRepo.Create(ctx, asset); err != nil {
		return nil, err
	}

	resp := asset.ToResponse()
	return &resp, nil
}

func (s *AssetService) UpdateAsset(ctx context.Context, id uuid.UUID, req domain.UpdateAssetRequest) (*domain.AssetResponse, error) {
	asset, err := s.assetRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrAssetNotFound
	}

	if req.Name != nil {
		asset.Name = *req.Name
	}
	if req.Status != nil {
		asset.Status = *req.Status
	}
	if req.Category != nil {
		asset.Category = *req.Category
	}
	if req.Location != nil {
		asset.Location = *req.Location
	}
	if req.PurchaseValue != nil {
		asset.PurchaseValue = *req.PurchaseValue
	}

	asset.UpdatedAt = time.Now()

	if err := s.assetRepo.Update(ctx, asset); err != nil {
		return nil, err
	}

	resp := asset.ToResponse()
	return &resp, nil
}

func (s *AssetService) DeleteAsset(ctx context.Context, id uuid.UUID) error {
	_, err := s.assetRepo.FindByID(ctx, id)
	if err != nil {
		return ErrAssetNotFound
	}
	return s.assetRepo.Delete(ctx, id)
}

func (s *AssetService) GetAsset(ctx context.Context, id uuid.UUID) (*domain.AssetResponse, error) {
	asset, err := s.assetRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrAssetNotFound
	}
	resp := asset.ToResponse()
	return &resp, nil
}

func (s *AssetService) GetAssetByCode(ctx context.Context, code string) (*domain.AssetResponse, error) {
	asset, err := s.assetRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, ErrAssetNotFound
	}
	resp := asset.ToResponse()
	return &resp, nil
}

func (s *AssetService) GetAllAssets(ctx context.Context) ([]domain.AssetResponse, error) {
	assets, err := s.assetRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.AssetResponse, len(assets))
	for i, a := range assets {
		responses[i] = a.ToResponse()
	}
	return responses, nil
}
