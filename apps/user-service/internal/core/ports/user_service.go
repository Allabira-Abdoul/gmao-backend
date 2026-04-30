package ports

import (
	"context"

	"backend-gmao/apps/user-service/internal/core/domain"
	"github.com/google/uuid"
)

// UserServicePort defines the primary port for user-related use cases.
type UserServicePort interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.InternalUserResponse, error)
	ListUsers(ctx context.Context, page, perPage int) ([]domain.UserResponse, int64, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// RoleServicePort defines the primary port for role-related use cases.
type RoleServicePort interface {
	CreateRole(ctx context.Context, req domain.CreateRoleRequest) (*domain.RoleResponse, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.RoleResponse, error)
	ListRoles(ctx context.Context) ([]domain.RoleResponse, error)
	UpdateRole(ctx context.Context, id uuid.UUID, req domain.UpdateRoleRequest) (*domain.RoleResponse, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error
	SetRolePrivileges(ctx context.Context, roleID uuid.UUID, req domain.SetPrivilegesRequest) (*domain.RoleResponse, error)
}
