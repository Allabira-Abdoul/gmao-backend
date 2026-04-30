package application

import (
	"context"
	"errors"
	"fmt"

	"backend-gmao/apps/user-service/internal/core/domain"
	"backend-gmao/apps/user-service/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrRoleNotFoundByID   = errors.New("role not found")
	ErrRoleLibelleExists  = errors.New("a role with this label already exists")
	ErrInvalidPrivileges  = errors.New("one or more privileges are invalid")
	ErrRoleHasUsers       = errors.New("cannot delete a role that is assigned to users")
)

// RoleService implements the RoleServicePort primary port.
type RoleService struct {
	roleRepo ports.RoleRepository
	userRepo ports.UserRepository
}

// NewRoleService creates a new RoleService instance.
func NewRoleService(roleRepo ports.RoleRepository, userRepo ports.UserRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// CreateRole creates a new role with validated privileges.
func (s *RoleService) CreateRole(ctx context.Context, req domain.CreateRoleRequest) (*domain.RoleResponse, error) {
	// Check if label already exists
	existing, _ := s.roleRepo.FindByLibelle(ctx, req.Libelle)
	if existing != nil {
		return nil, ErrRoleLibelleExists
	}

	// Validate all privileges are system-defined
	invalidPrivs := domain.ValidatePrivileges(req.Privileges)
	if len(invalidPrivs) > 0 {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPrivileges, invalidPrivs)
	}

	// Build the role with privileges
	role := &domain.Role{
		Libelle:     req.Libelle,
		Description: req.Description,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Set privileges
	if err := s.roleRepo.SetPrivileges(ctx, role.IDRole, req.Privileges); err != nil {
		return nil, fmt.Errorf("failed to set privileges: %w", err)
	}

	// Reload with privileges
	created, err := s.roleRepo.FindByID(ctx, role.IDRole)
	if err != nil {
		return nil, fmt.Errorf("failed to reload role: %w", err)
	}

	resp := created.ToResponse()
	return &resp, nil
}

// GetRoleByID retrieves a role by its UUID.
func (s *RoleService) GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoleNotFoundByID
	}

	resp := role.ToResponse()
	return &resp, nil
}

// ListRoles returns all roles.
func (s *RoleService) ListRoles(ctx context.Context) ([]domain.RoleResponse, error) {
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	responses := make([]domain.RoleResponse, 0, len(roles))
	for _, r := range roles {
		responses = append(responses, r.ToResponse())
	}

	return responses, nil
}

// UpdateRole updates an existing role's label and/or description.
func (s *RoleService) UpdateRole(ctx context.Context, id uuid.UUID, req domain.UpdateRoleRequest) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoleNotFoundByID
	}

	if req.Libelle != nil {
		// Check uniqueness
		existing, _ := s.roleRepo.FindByLibelle(ctx, *req.Libelle)
		if existing != nil && existing.IDRole != id {
			return nil, ErrRoleLibelleExists
		}
		role.Libelle = *req.Libelle
	}

	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Reload
	updated, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload role: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// DeleteRole removes a role if it has no assigned users.
func (s *RoleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	_, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return ErrRoleNotFoundByID
	}

	// TODO: Check if any users are assigned to this role before deleting.
	// For now, the database FK constraint will prevent deletion if users exist.

	if err := s.roleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// SetRolePrivileges replaces all privileges for a role.
func (s *RoleService) SetRolePrivileges(ctx context.Context, roleID uuid.UUID, req domain.SetPrivilegesRequest) (*domain.RoleResponse, error) {
	_, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, ErrRoleNotFoundByID
	}

	// Validate all privileges
	invalidPrivs := domain.ValidatePrivileges(req.Privileges)
	if len(invalidPrivs) > 0 {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPrivileges, invalidPrivs)
	}

	if err := s.roleRepo.SetPrivileges(ctx, roleID, req.Privileges); err != nil {
		return nil, fmt.Errorf("failed to set privileges: %w", err)
	}

	// Reload
	updated, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload role: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}
