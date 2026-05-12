package postgres

import (
	"context"
	"fmt"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InterventionRepository is the GORM-based implementation of the InterventionRepository port.
type InterventionRepository struct {
	db *gorm.DB
}

// NewInterventionRepository creates a new InterventionRepository instance.
func NewInterventionRepository(db *gorm.DB) *InterventionRepository {
	return &InterventionRepository{db: db}
}

// Create persists a new intervention to the database.
func (r *InterventionRepository) Create(ctx context.Context, intervention *domain.Intervention) error {
	result := r.db.WithContext(ctx).Create(intervention)
	if result.Error != nil {
		return fmt.Errorf("postgres create intervention: %w", result.Error)
	}
	return nil
}

// GetByID retrieves an intervention by UUID.
func (r *InterventionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Intervention, error) {
	var intervention domain.Intervention
	result := r.db.WithContext(ctx).
		Where("id_intervention = ?", id).
		First(&intervention)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find intervention by id: %w", result.Error)
	}
	return &intervention, nil
}

// ListByOrdreTravail retrieves interventions for a specific work order.
func (r *InterventionRepository) ListByOrdreTravail(ctx context.Context, ordreTravailID uuid.UUID, limit, offset int) ([]domain.Intervention, int64, error) {
	var interventions []domain.Intervention
	var total int64

	r.db.WithContext(ctx).Model(&domain.Intervention{}).Where("id_ordre_travail = ?", ordreTravailID).Count(&total)

	result := r.db.WithContext(ctx).
		Where("id_ordre_travail = ?", ordreTravailID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&interventions)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("postgres list interventions by ordre_travail: %w", result.Error)
	}
	return interventions, total, nil
}

// ListByTechnicien retrieves interventions for a specific technician.
func (r *InterventionRepository) ListByTechnicien(ctx context.Context, technicienID uuid.UUID, limit, offset int) ([]domain.Intervention, int64, error) {
	var interventions []domain.Intervention
	var total int64

	r.db.WithContext(ctx).Model(&domain.Intervention{}).Where("id_technicien = ?", technicienID).Count(&total)

	result := r.db.WithContext(ctx).
		Where("id_technicien = ?", technicienID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&interventions)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("postgres list interventions by technicien: %w", result.Error)
	}
	return interventions, total, nil
}

// Update updates an existing intervention in the database.
func (r *InterventionRepository) Update(ctx context.Context, intervention *domain.Intervention) error {
	result := r.db.WithContext(ctx).Save(intervention)
	if result.Error != nil {
		return fmt.Errorf("postgres update intervention: %w", result.Error)
	}
	return nil
}

// Delete removes an intervention from the database by UUID.
func (r *InterventionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id_intervention = ?", id).Delete(&domain.Intervention{})
	if result.Error != nil {
		return fmt.Errorf("postgres delete intervention: %w", result.Error)
	}
	return nil
}
