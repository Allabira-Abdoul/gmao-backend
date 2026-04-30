package application

import (
	"context"
	"errors"
	"fmt"

	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/apps/user-service/internal/core/ports"
	"backend-gmao/pkg/auth"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailExists      = errors.New("a user with this email already exists")
	ErrRoleNotFound     = errors.New("the specified role does not exist")
	ErrInvalidAccount   = errors.New("user account is not active")
	ErrCannotDeleteSelf = errors.New("you cannot delete your own account")
)

// UserService implements the UserServicePort primary port.
type UserService struct {
	userRepo ports.UserRepository
	roleRepo ports.RoleRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepo ports.UserRepository, roleRepo ports.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateUser creates a new user after validating business rules.
func (s *UserService) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
	// Check if email already exists
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailExists
	}

	// Validate role exists
	roleID, err := uuid.Parse(req.IDRole)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID: %w", err)
	}

	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil || role == nil {
		return nil, ErrRoleNotFound
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(req.MotDePasse)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		NomComplet:   req.NomComplet,
		Email:        req.Email,
		MotDePasse:   hashedPassword,
		StatutCompte: domain.StatusActive,
		IDRole:       roleID,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Reload the user with role preloaded
	created, err := s.userRepo.FindByID(ctx, user.IDUtilisateur)
	if err != nil {
		return nil, fmt.Errorf("failed to reload user: %w", err)
	}

	resp := created.ToResponse()
	return &resp, nil
}

// GetUserByID retrieves a user by their UUID.
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	resp := user.ToResponse()
	return &resp, nil
}

// GetUserByEmail retrieves a user by email for internal authentication use.
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.InternalUserResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	resp := user.ToInternalResponse()
	return &resp, nil
}

// ListUsers returns a paginated list of users.
func (s *UserService) ListUsers(ctx context.Context, page, perPage int) ([]domain.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	users, total, err := s.userRepo.FindAll(ctx, offset, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	responses := make([]domain.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}

	return responses, total, nil
}

// UpdateUser updates an existing user's fields.
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.NomComplet != nil {
		user.NomComplet = *req.NomComplet
	}

	if req.Email != nil {
		// Check email uniqueness
		existing, _ := s.userRepo.FindByEmail(ctx, *req.Email)
		if existing != nil && existing.IDUtilisateur != id {
			return nil, ErrEmailExists
		}
		user.Email = *req.Email
	}

	if req.StatutCompte != nil {
		user.StatutCompte = domain.AccountStatus(*req.StatutCompte)
	}

	if req.IDRole != nil {
		roleID, err := uuid.Parse(*req.IDRole)
		if err != nil {
			return nil, fmt.Errorf("invalid role ID: %w", err)
		}
		role, err := s.roleRepo.FindByID(ctx, roleID)
		if err != nil || role == nil {
			return nil, ErrRoleNotFound
		}
		user.IDRole = roleID
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Reload with role
	updated, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload user: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// DeleteUser removes a user by their UUID.
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
