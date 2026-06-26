package secondary

import (
	"context"
	"errors"

	"backend-gmao/apps/auth-service/internal/core/domain"
)

var ErrUserNotFound = errors.New("user not found")

// UserClient defines the interface for communicating with the user-service.
type UserClient interface {
	// GetUserByEmail retrieves a user by their email.
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
