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
	ErrOrdreTravailNotFound = errors.New("work order not found")
	ErrReferenceExists      = errors.New("a work order with this reference already exists")
)

// OrdreTravailService implements the OrdreTravailService port.
type OrdreTravailService struct {
	repo ports.OrdreTravailRepository
}

// NewOrdreTravailService creates a new OrdreTravailService instance.
func NewOrdreTravailService(repo ports.OrdreTravailRepository) *OrdreTravailService {
	return &OrdreTravailService{repo: repo}
}

// generateReference creates a unique reference for a work order (OT-YYYYMMDD-XXXX).
func generateReference() string {
	now := time.Now()
	short := uuid.New().String()[:4]
	return fmt.Sprintf("OT-%s-%s", now.Format("20060102"), short)
}

// CreateOrdreTravail creates a new work order.
func (s *OrdreTravailService) CreateOrdreTravail(ctx context.Context, req *domain.CreateOrdreTravailRequest, createdBy uuid.UUID) (*domain.OrdreTravailResponse, error) {
	equipementID, err := uuid.Parse(req.IDEquipement)
	if err != nil {
		return nil, fmt.Errorf("invalid equipement ID: %w", err)
	}

	priorite := domain.Priorite(req.Priorite)
	if req.Priorite == "" {
		priorite = domain.PrioriteMoyenne
	}

	ordre := &domain.OrdreTravail{
		Reference:         generateReference(),
		Titre:             req.Titre,
		Description:       req.Description,
		Statut:            domain.StatutPlanifie,
		TypeMaintenance:   domain.TypeMaintenance(req.TypeMaintenance),
		Priorite:          priorite,
		IDEquipement:      equipementID,
		IDUtilisateurCree: createdBy,
	}

	if req.IDUtilisateurAssigne != nil {
		assigneID, err := uuid.Parse(*req.IDUtilisateurAssigne)
		if err != nil {
			return nil, fmt.Errorf("invalid assigned user ID: %w", err)
		}
		ordre.IDUtilisateurAssigne = &assigneID
	}

	if req.DateDebutPrevue != nil {
		t, err := time.Parse(time.RFC3339, *req.DateDebutPrevue)
		if err != nil {
			return nil, fmt.Errorf("invalid date_debut_prevue format: %w", err)
		}
		ordre.DateDebutPrevue = &t
	}

	if req.DateFinPrevue != nil {
		t, err := time.Parse(time.RFC3339, *req.DateFinPrevue)
		if err != nil {
			return nil, fmt.Errorf("invalid date_fin_prevue format: %w", err)
		}
		ordre.DateFinPrevue = &t
	}

	if req.Commentaire != "" {
		ordre.Commentaire = req.Commentaire
	}

	if err := s.repo.Create(ctx, ordre); err != nil {
		return nil, fmt.Errorf("failed to create work order: %w", err)
	}

	created, err := s.repo.GetByID(ctx, ordre.IDOrdreTravail)
	if err != nil {
		return nil, fmt.Errorf("failed to reload work order: %w", err)
	}

	resp := created.ToResponse()
	return &resp, nil
}

// GetOrdreTravailByID retrieves a work order by its UUID.
func (s *OrdreTravailService) GetOrdreTravailByID(ctx context.Context, id uuid.UUID) (*domain.OrdreTravailResponse, error) {
	ordre, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrOrdreTravailNotFound
	}
	resp := ordre.ToResponse()
	return &resp, nil
}

// ListOrdresTravail returns a paginated list of work orders.
func (s *OrdreTravailService) ListOrdresTravail(ctx context.Context, page, perPage int) ([]domain.OrdreTravailResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	ordres, total, err := s.repo.List(ctx, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list work orders: %w", err)
	}

	responses := make([]domain.OrdreTravailResponse, 0, len(ordres))
	for _, o := range ordres {
		responses = append(responses, o.ToResponse())
	}
	return responses, total, nil
}

// ListByEquipement returns work orders for a specific equipment.
func (s *OrdreTravailService) ListByEquipement(ctx context.Context, equipementID uuid.UUID, page, perPage int) ([]domain.OrdreTravailResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	ordres, total, err := s.repo.ListByEquipement(ctx, equipementID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list work orders by equipement: %w", err)
	}

	responses := make([]domain.OrdreTravailResponse, 0, len(ordres))
	for _, o := range ordres {
		responses = append(responses, o.ToResponse())
	}
	return responses, total, nil
}

// ListByStatut returns work orders by status.
func (s *OrdreTravailService) ListByStatut(ctx context.Context, statut domain.OrdreTravailStatut, page, perPage int) ([]domain.OrdreTravailResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	ordres, total, err := s.repo.ListByStatut(ctx, statut, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list work orders by statut: %w", err)
	}

	responses := make([]domain.OrdreTravailResponse, 0, len(ordres))
	for _, o := range ordres {
		responses = append(responses, o.ToResponse())
	}
	return responses, total, nil
}

// ListByAssigne returns work orders assigned to a specific user.
func (s *OrdreTravailService) ListByAssigne(ctx context.Context, utilisateurID uuid.UUID, page, perPage int) ([]domain.OrdreTravailResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	ordres, total, err := s.repo.ListByAssigne(ctx, utilisateurID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list work orders by assigne: %w", err)
	}

	responses := make([]domain.OrdreTravailResponse, 0, len(ordres))
	for _, o := range ordres {
		responses = append(responses, o.ToResponse())
	}
	return responses, total, nil
}

// UpdateOrdreTravail updates an existing work order's fields.
func (s *OrdreTravailService) UpdateOrdreTravail(ctx context.Context, id uuid.UUID, req *domain.UpdateOrdreTravailRequest) (*domain.OrdreTravailResponse, error) {
	ordre, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrOrdreTravailNotFound
	}

	if req.Titre != nil {
		ordre.Titre = *req.Titre
	}
	if req.Description != nil {
		ordre.Description = *req.Description
	}
	if req.Statut != nil {
		ordre.Statut = domain.OrdreTravailStatut(*req.Statut)
	}
	if req.Priorite != nil {
		ordre.Priorite = domain.Priorite(*req.Priorite)
	}
	if req.IDUtilisateurAssigne != nil {
		assigneID, err := uuid.Parse(*req.IDUtilisateurAssigne)
		if err != nil {
			return nil, fmt.Errorf("invalid assigned user ID: %w", err)
		}
		ordre.IDUtilisateurAssigne = &assigneID
	}
	if req.DateDebutPrevue != nil {
		t, err := time.Parse(time.RFC3339, *req.DateDebutPrevue)
		if err != nil {
			return nil, fmt.Errorf("invalid date_debut_prevue: %w", err)
		}
		ordre.DateDebutPrevue = &t
	}
	if req.DateFinPrevue != nil {
		t, err := time.Parse(time.RFC3339, *req.DateFinPrevue)
		if err != nil {
			return nil, fmt.Errorf("invalid date_fin_prevue: %w", err)
		}
		ordre.DateFinPrevue = &t
	}
	if req.DateDebutReelle != nil {
		t, err := time.Parse(time.RFC3339, *req.DateDebutReelle)
		if err != nil {
			return nil, fmt.Errorf("invalid date_debut_reelle: %w", err)
		}
		ordre.DateDebutReelle = &t
	}
	if req.DateFinReelle != nil {
		t, err := time.Parse(time.RFC3339, *req.DateFinReelle)
		if err != nil {
			return nil, fmt.Errorf("invalid date_fin_reelle: %w", err)
		}
		ordre.DateFinReelle = &t
	}
	if req.Commentaire != nil {
		ordre.Commentaire = *req.Commentaire
	}

	if err := s.repo.Update(ctx, ordre); err != nil {
		return nil, fmt.Errorf("failed to update work order: %w", err)
	}

	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload work order: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// DeleteOrdreTravail removes a work order by its UUID.
func (s *OrdreTravailService) DeleteOrdreTravail(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrOrdreTravailNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete work order: %w", err)
	}
	return nil
}
