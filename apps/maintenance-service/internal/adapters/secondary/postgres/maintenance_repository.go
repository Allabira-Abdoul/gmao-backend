package postgres

import (
	"context"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type maintenanceRepository struct {
	db *gorm.DB
}

// NewMaintenanceRepository creates a GORM maintenance repository.
func NewMaintenanceRepository(db *gorm.DB) *maintenanceRepository {
	return &maintenanceRepository{db: db}
}

func (r *maintenanceRepository) CreateWorkOrder(ctx context.Context, wo *domain.OrdreTravail) error {
	return r.db.WithContext(ctx).Create(wo).Error
}

func (r *maintenanceRepository) UpdateWorkOrder(ctx context.Context, wo *domain.OrdreTravail) error {
	return r.db.WithContext(ctx).Save(wo).Error
}

func (r *maintenanceRepository) DeleteWorkOrder(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.OrdreTravail{}, "id = ?", id).Error
}

func (r *maintenanceRepository) FindWorkOrderByID(ctx context.Context, id uuid.UUID) (*domain.OrdreTravail, error) {
	var wo domain.OrdreTravail
	if err := r.db.WithContext(ctx).First(&wo, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wo, nil
}

func (r *maintenanceRepository) FindAllWorkOrders(ctx context.Context) ([]domain.OrdreTravail, error) {
	var wos []domain.OrdreTravail
	// ⚡ Bolt Optimization: Use Preload to eagerly fetch interventions and their measurements
	// in a single query to eliminate N+1 queries when fetching all work orders.
	if err := r.db.WithContext(ctx).Preload("Interventions").Preload("Interventions.Measurements").Find(&wos).Error; err != nil {
		return nil, err
	}
	return wos, nil
}

func (r *maintenanceRepository) CreateIntervention(ctx context.Context, intervention *domain.Intervention) error {
	return r.db.WithContext(ctx).Create(intervention).Error
}

func (r *maintenanceRepository) FindInterventionsByWorkOrderID(ctx context.Context, workOrderID uuid.UUID) ([]domain.Intervention, error) {
	var interventions []domain.Intervention
	if err := r.db.WithContext(ctx).Preload("Measurements").Where("work_order_id = ?", workOrderID).Find(&interventions).Error; err != nil {
		return nil, err
	}
	return interventions, nil
}
