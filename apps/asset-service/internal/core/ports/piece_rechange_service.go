package ports

import (
	"context"

	"github.com/google/uuid"

	"backend-gmao/apps/asset-service/internal/core/domain"
)

// PieceRechangeService defines the use cases for PieceRechange operations.
type PieceRechangeService interface {
	CreatePieceRechange(ctx context.Context, req *domain.CreatePieceRechangeRequest) (*domain.PieceRechangeResponse, error)
	GetPieceRechangeByID(ctx context.Context, id uuid.UUID) (*domain.PieceRechangeResponse, error)
	GetPieceRechangeByReference(ctx context.Context, reference string) (*domain.PieceRechangeResponse, error)
	ListPiecesRechange(ctx context.Context, limit, offset int) ([]domain.PieceRechangeResponse, int64, error)
	UpdatePieceRechange(ctx context.Context, id uuid.UUID, req *domain.UpdatePieceRechangeRequest) (*domain.PieceRechangeResponse, error)
	DeletePieceRechange(ctx context.Context, id uuid.UUID) error
}
