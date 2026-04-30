package postgres

import (
	"context"
	"fmt"

	"backend-gmao/apps/user-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository is the GORM-based implementation of the UserRepository port.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create persists a new user to the database.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("postgres create user: %w", result.Error)
	}
	return nil
}

// FindByID retrieves a user by UUID, preloading their role and role privileges.
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Role.Privileges").
		Where("id_utilisateur = ?", id).
		First(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find user by id: %w", result.Error)
	}
	return &user, nil
}

// FindByEmail retrieves a user by email, preloading their role and role privileges.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Role.Privileges").
		Where("email = ?", email).
		First(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find user by email: %w", result.Error)
	}
	return &user, nil
}

// FindAll retrieves a paginated list of users with their roles.
func (r *UserRepository) FindAll(ctx context.Context, offset, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	// Count total
	r.db.WithContext(ctx).Model(&domain.User{}).Count(&total)

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Role.Privileges").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("postgres find all users: %w", result.Error)
	}

	return users, total, nil
}

// Update updates an existing user in the database.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("postgres update user: %w", result.Error)
	}
	return nil
}

// Delete removes a user from the database by UUID.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id_utilisateur = ?", id).Delete(&domain.User{})
	if result.Error != nil {
		return fmt.Errorf("postgres delete user: %w", result.Error)
	}
	return nil
}
