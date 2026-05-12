package postgres

import (
	"context"
	"fmt"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EquipementRepository is the GORM-based implementation of the EquipementRepository port.
type EquipementRepository struct {
	db *gorm.DB
}

// NewEquipementRepository creates a new EquipementRepository instance.
func NewEquipementRepository(db *gorm.DB) *EquipementRepository {
	return &EquipementRepository{db: db}
}

// Create persists a new equipement to the database.
func (r *EquipementRepository) Create(ctx context.Context, equipement *domain.Equipement) error {
	result := r.db.WithContext(ctx).Create(equipement)
	if result.Error != nil {
		return fmt.Errorf("postgres create equipement: %w", result.Error)
	}
	return nil
}

// GetByID retrieves an equipement by UUID.
func (r *EquipementRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Equipement, error) {
	var equipement domain.Equipement
	result := r.db.WithContext(ctx).
		Where("id_equipement = ?", id).
		First(&equipement)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find equipement by id: %w", result.Error)
	}
	return &equipement, nil
}

// GetByCode retrieves an equipement by its unique code.
func (r *EquipementRepository) GetByCode(ctx context.Context, code string) (*domain.Equipement, error) {
	var equipement domain.Equipement
	result := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&equipement)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find equipement by code: %w", result.Error)
	}
	return &equipement, nil
}

// List retrieves a paginated list of equipements.
func (r *EquipementRepository) List(ctx context.Context, limit, offset int) ([]domain.Equipement, int64, error) {
	var equipements []domain.Equipement
	var total int64

	r.db.WithContext(ctx).Model(&domain.Equipement{}).Count(&total)

	result := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&equipements)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("postgres find all equipements: %w", result.Error)
	}

	return equipements, total, nil
}

// Update updates an existing equipement in the database.
func (r *EquipementRepository) Update(ctx context.Context, equipement *domain.Equipement) error {
	result := r.db.WithContext(ctx).Save(equipement)
	if result.Error != nil {
		return fmt.Errorf("postgres update equipement: %w", result.Error)
	}
	return nil
}

// Delete removes an equipement from the database by UUID.
func (r *EquipementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id_equipement = ?", id).Delete(&domain.Equipement{})
	if result.Error != nil {
		return fmt.Errorf("postgres delete equipement: %w", result.Error)
	}
	return nil
}
