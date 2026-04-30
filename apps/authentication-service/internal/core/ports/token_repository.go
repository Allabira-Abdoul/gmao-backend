package ports

import (
	"context"

	"backend-gmao/apps/authentication-service/internal/core/domain"
	"github.com/google/uuid"
)

// TokenRepository defines the secondary port for refresh token persistence.
type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredTokens(ctx context.Context) error
}
