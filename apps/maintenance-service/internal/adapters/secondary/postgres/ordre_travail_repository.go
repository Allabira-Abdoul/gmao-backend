package postgres

import (
	"context"
	"fmt"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrdreTravailRepository is the GORM-based implementation of the OrdreTravailRepository port.
type OrdreTravailRepository struct {
	db *gorm.DB
}

// NewOrdreTravailRepository creates a new OrdreTravailRepository instance.
func NewOrdreTravailRepository(db *gorm.DB) *OrdreTravailRepository {
	return &OrdreTravailRepository{db: db}
}

// Create persists a new work order to the database.
func (r *OrdreTravailRepository) Create(ctx context.Context, ordre *domain.OrdreTravail) error {
	result := r.db.WithContext(ctx).Create(ordre)
	if result.Error != nil {
		return fmt.Errorf("postgres create ordre_travail: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a work order by UUID.
func (r *OrdreTravailRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.OrdreTravail, error) {
	var ordre domain.OrdreTravail
	result := r.db.WithContext(ctx).
		Where("id_ordre_travail = ?", id).
		First(&ordre)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find ordre_travail by id: %w", result.Error)
	}
	return &ordre, nil
}

// GetByReference retrieves a work order by its reference code.
func (r *OrdreTravailRepository) GetByReference(ctx context.Context, reference string) (*domain.OrdreTravail, error) {
	var ordre domain.OrdreTravail
	result := r.db.WithContext(ctx).
		Where("reference = ?", reference).
		First(&ordre)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres find ordre_travail by reference: %w", result.Error)
	}
	return &ordre, nil
}

// List retrieves a paginated list of work orders.
func (r *OrdreTravailRepository) List(ctx context.Context, limit, offset int) ([]domain.OrdreTravail, int64, error) {
	var ordres []domain.OrdreTravail
	var total int64

	r.db.WithContext(ctx).Model(&domain.OrdreTravail{}).Count(&total)

	result := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ordres)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("postgres list ordres_travail: %w", result.Error)
	}
	return ordres, total, nil
}

// ListByEquipement retrieves work orders for a specific equipment.
func (r *OrdreTravailRepository) ListByEquipement(ctx context.Context, equipementID uuid.UUID, limit, offset int) ([]domain.OrdreTravail, int64, error) {
	var ordres []domain.OrdreTravail
	var total int64

	r.db.WithContext(ctx).Model(&domain.OrdreTravail{}).Where("id_equipement = ?", equipementID).Count(&total)

	result := r.db.WithContext(ctx).
		Where("id_equipement = ?", equipementID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ordres)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("postgres list ordres_travail by equipement: %w", result.Error)
	}
	return ordres, total, nil
}

// ListByStatut retrieves work orders by status.
func (r *OrdreTravailRepository) ListByStatut(ctx context.Context, statut domain.OrdreTravailStatut, limit, offset int) ([]domain.OrdreTravail, int64, error) {
	var ordres []domain.OrdreTravail
	var total int64

	r.db.WithContext(ctx).Model(&domain.OrdreTravail{}).Where("statut = ?", statut).Count(&total)

	result := r.db.WithContext(ctx).
		Where("statut = ?", statut).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ordres)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("postgres list ordres_travail by statut: %w", result.Error)
	}
	return ordres, total, nil
}

// ListByAssigne retrieves work orders assigned to a specific user.
func (r *OrdreTravailRepository) ListByAssigne(ctx context.Context, utilisateurID uuid.UUID, limit, offset int) ([]domain.OrdreTravail, int64, error) {
	var ordres []domain.OrdreTravail
	var total int64

	r.db.WithContext(ctx).Model(&domain.OrdreTravail{}).Where("id_utilisateur_assigne = ?", utilisateurID).Count(&total)

	result := r.db.WithContext(ctx).
		Where("id_utilisateur_assigne = ?", utilisateurID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&ordres)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("postgres list ordres_travail by assigne: %w", result.Error)
	}
	return ordres, total, nil
}

// Update updates an existing work order in the database.
func (r *OrdreTravailRepository) Update(ctx context.Context, ordre *domain.OrdreTravail) error {
	result := r.db.WithContext(ctx).Save(ordre)
	if result.Error != nil {
		return fmt.Errorf("postgres update ordre_travail: %w", result.Error)
	}
	return nil
}

// Delete removes a work order from the database by UUID.
func (r *OrdreTravailRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id_ordre_travail = ?", id).Delete(&domain.OrdreTravail{})
	if result.Error != nil {
		return fmt.Errorf("postgres delete ordre_travail: %w", result.Error)
	}
	return nil
}
