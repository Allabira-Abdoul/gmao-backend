package ports

import (
	"context"

	"github.com/google/uuid"

	"backend-gmao/apps/asset-service/internal/core/domain"
)

// PieceRechangeRepository defines the interface for data access operations related to PieceRechange.
type PieceRechangeRepository interface {
	Create(ctx context.Context, piece *domain.PieceRechange) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PieceRechange, error)
	GetByReference(ctx context.Context, reference string) (*domain.PieceRechange, error)
	List(ctx context.Context, limit, offset int) ([]domain.PieceRechange, int64, error)
	Update(ctx context.Context, piece *domain.PieceRechange) error
	Delete(ctx context.Context, id uuid.UUID) error
}
