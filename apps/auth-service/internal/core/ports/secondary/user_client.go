package secondary

import (
	"context"

	"backend-gmao/apps/auth-service/internal/core/domain"
)

// UserClient defines the interface for communicating with the user service.
// This is an application of the Dependency Inversion Principle (DIP),
// ensuring the application service depends on abstractions rather than
// concrete HTTP clients and service registries.
type UserClient interface {
	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
