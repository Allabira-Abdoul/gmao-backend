package ports

import (
	"context"

	"backend-gmao/apps/authentication-service/internal/core/domain"
)

// AuthServicePort defines the primary port for authentication use cases.
type AuthServicePort interface {
	Login(ctx context.Context, req domain.LoginRequest) (*domain.TokenPair, error)
	RefreshToken(ctx context.Context, req domain.RefreshRequest) (*domain.TokenPair, error)
	Logout(ctx context.Context, req domain.LogoutRequest) error
}
