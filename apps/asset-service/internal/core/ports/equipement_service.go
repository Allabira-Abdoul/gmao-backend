package ports

import (
	"context"

	"github.com/google/uuid"

	"backend-gmao/apps/asset-service/internal/core/domain"
)

// EquipementService defines the use cases for Equipement operations.
type EquipementService interface {
	CreateEquipement(ctx context.Context, req *domain.CreateEquipementRequest) (*domain.EquipementResponse, error)
	GetEquipementByID(ctx context.Context, id uuid.UUID) (*domain.EquipementResponse, error)
	GetEquipementByCode(ctx context.Context, code string) (*domain.EquipementResponse, error)
	ListEquipements(ctx context.Context, limit, offset int) ([]domain.EquipementResponse, int64, error)
	UpdateEquipement(ctx context.Context, id uuid.UUID, req *domain.UpdateEquipementRequest) (*domain.EquipementResponse, error)
	DeleteEquipement(ctx context.Context, id uuid.UUID) error
}
