package secondary

import (
	"context"
	"errors"
	"backend-gmao/apps/auth-service/internal/core/domain"
)

var ErrUserNotFound = errors.New("user not found")

// UserClient defines the port to interact with the user service.
// 🏛️ Atlas: Dependency Inversion Principle (DIP) applied here.
// Abstracting the user fetching logic into a dedicated port decouple the
// application service layer from low-level network and discovery operations.
type UserClient interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
