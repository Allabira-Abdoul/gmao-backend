package ports

import (
	"context"

	"backend-gmao/apps/authentication-service/internal/core/domain"

	"github.com/google/uuid"
)

// UserClient defines the secondary port for communicating with the user-service.
type UserClient interface {
	FindUserByEmail(ctx context.Context, email string) (*domain.UserInfo, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*domain.UserInfo, error)
}
