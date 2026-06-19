package secondary

import (
	"context"

	"backend-gmao/apps/auth-service/internal/core/domain"
)

// UserClient defines the interface for communicating with the user-service (Secondary Port)
type UserClient interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
