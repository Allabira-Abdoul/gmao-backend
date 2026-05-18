package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/apps/maintenance-service/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrInterventionNotFound = errors.New("intervention not found")
)

// InterventionService implements the InterventionService port.
type InterventionService struct {
	repo ports.InterventionRepository
}

// NewInterventionService creates a new InterventionService instance.
func NewInterventionService(repo ports.InterventionRepository) *InterventionService {
	return &InterventionService{repo: repo}
}

// CreateIntervention creates a new intervention.
func (s *InterventionService) CreateIntervention(ctx context.Context, req *domain.CreateInterventionRequest) (*domain.InterventionResponse, error) {
	ordreTravailID, err := uuid.Parse(req.IDOrdreTravail)
	if err != nil {
		return nil, fmt.Errorf("invalid work order ID: %w", err)
	}

	technicienID, err := uuid.Parse(req.IDTechnicien)
	if err != nil {
		return nil, fmt.Errorf("invalid technician ID: %w", err)
	}

	dateDebut, err := time.Parse(time.RFC3339, req.DateDebut)
	if err != nil {
		return nil, fmt.Errorf("invalid date_debut format: %w", err)
	}

	intervention := &domain.Intervention{
		IDOrdreTravail: ordreTravailID,
		IDTechnicien:   technicienID,
		Statut:         domain.InterventionEnCours,
		DateDebut:      dateDebut,
	}

	if err := s.repo.Create(ctx, intervention); err != nil {
		return nil, fmt.Errorf("failed to create intervention: %w", err)
	}

	created, err := s.repo.GetByID(ctx, intervention.IDIntervention)
	if err != nil {
		return nil, fmt.Errorf("failed to reload intervention: %w", err)
	}

	resp := created.ToResponse()
	return &resp, nil
}

// GetInterventionByID retrieves an intervention by its UUID.
func (s *InterventionService) GetInterventionByID(ctx context.Context, id uuid.UUID) (*domain.InterventionResponse, error) {
	intervention, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrInterventionNotFound
	}
	resp := intervention.ToResponse()
	return &resp, nil
}

// ListByOrdreTravail returns interventions for a specific work order.
func (s *InterventionService) ListByOrdreTravail(ctx context.Context, ordreTravailID uuid.UUID, limit, offset int) ([]domain.InterventionResponse, int64, error) {
	interventions, total, err := s.repo.ListByOrdreTravail(ctx, ordreTravailID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list interventions: %w", err)
	}

	responses := make([]domain.InterventionResponse, 0, len(interventions))
	for _, i := range interventions {
		responses = append(responses, i.ToResponse())
	}
	return responses, total, nil
}

// ListByTechnicien returns interventions for a specific technician.
func (s *InterventionService) ListByTechnicien(ctx context.Context, technicienID uuid.UUID, limit, offset int) ([]domain.InterventionResponse, int64, error) {
	interventions, total, err := s.repo.ListByTechnicien(ctx, technicienID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list interventions by technicien: %w", err)
	}

	responses := make([]domain.InterventionResponse, 0, len(interventions))
	for _, i := range interventions {
		responses = append(responses, i.ToResponse())
	}
	return responses, total, nil
}

// UpdateIntervention updates an existing intervention.
func (s *InterventionService) UpdateIntervention(ctx context.Context, id uuid.UUID, req *domain.UpdateInterventionRequest) (*domain.InterventionResponse, error) {
	intervention, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrInterventionNotFound
	}

	if req.Statut != nil {
		intervention.Statut = domain.InterventionStatut(*req.Statut)
	}
	if req.DateFin != nil {
		t, err := time.Parse(time.RFC3339, *req.DateFin)
		if err != nil {
			return nil, fmt.Errorf("invalid date_fin format: %w", err)
		}
		intervention.DateFin = &t
	}
	if req.DureeMinutes != nil {
		intervention.DureeMinutes = req.DureeMinutes
	}
	if req.RapportIntervention != nil {
		intervention.RapportIntervention = *req.RapportIntervention
	}
	if req.ActionsEffectuees != nil {
		intervention.ActionsEffectuees = *req.ActionsEffectuees
	}

	if err := s.repo.Update(ctx, intervention); err != nil {
		return nil, fmt.Errorf("failed to update intervention: %w", err)
	}

	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload intervention: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// DeleteIntervention removes an intervention by its UUID.
func (s *InterventionService) DeleteIntervention(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrInterventionNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete intervention: %w", err)
	}
	return nil
}
