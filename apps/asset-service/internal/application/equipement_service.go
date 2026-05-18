package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"backend-gmao/apps/asset-service/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrEquipementNotFound = errors.New("equipement not found")
	ErrCodeExists         = errors.New("an equipement with this code already exists")
)

// EquipementService implements the EquipementService port.
type EquipementService struct {
	repo ports.EquipementRepository
}

// NewEquipementService creates a new EquipementService instance.
func NewEquipementService(repo ports.EquipementRepository) *EquipementService {
	return &EquipementService{repo: repo}
}

// CreateEquipement creates a new equipement after validating business rules.
func (s *EquipementService) CreateEquipement(ctx context.Context, req *domain.CreateEquipementRequest) (*domain.EquipementResponse, error) {
	// Check code uniqueness
	existing, _ := s.repo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, ErrCodeExists
	}

	dateAcq, err := time.Parse("2006-01-02", req.DateAcquisition)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	criticite := req.Criticite
	if criticite == "" {
		criticite = "BASSE"
	}

	equipement := &domain.Equipement{
		Nom:             req.Nom,
		Code:            req.Code,
		Description:     req.Description,
		Statut:          domain.StatusEnService,
		DateAcquisition: dateAcq,
		Localisation:    req.Localisation,
		Criticite:       criticite,
	}

	if err := s.repo.Create(ctx, equipement); err != nil {
		return nil, fmt.Errorf("failed to create equipement: %w", err)
	}

	created, err := s.repo.GetByID(ctx, equipement.IDEquipement)
	if err != nil {
		return nil, fmt.Errorf("failed to reload equipement: %w", err)
	}

	resp := created.ToResponse()
	return &resp, nil
}

// GetEquipementByID retrieves an equipement by its UUID.
func (s *EquipementService) GetEquipementByID(ctx context.Context, id uuid.UUID) (*domain.EquipementResponse, error) {
	equipement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrEquipementNotFound
	}

	resp := equipement.ToResponse()
	return &resp, nil
}

// GetEquipementByCode retrieves an equipement by its unique code.
func (s *EquipementService) GetEquipementByCode(ctx context.Context, code string) (*domain.EquipementResponse, error) {
	equipement, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, ErrEquipementNotFound
	}

	resp := equipement.ToResponse()
	return &resp, nil
}

// ListEquipements returns a paginated list of equipements.
func (s *EquipementService) ListEquipements(ctx context.Context, limit, offset int) ([]domain.EquipementResponse, int64, error) {
	equipements, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list equipements: %w", err)
	}

	responses := make([]domain.EquipementResponse, 0, len(equipements))
	for _, e := range equipements {
		responses = append(responses, e.ToResponse())
	}

	return responses, total, nil
}

// UpdateEquipement updates an existing equipement's fields.
func (s *EquipementService) UpdateEquipement(ctx context.Context, id uuid.UUID, req *domain.UpdateEquipementRequest) (*domain.EquipementResponse, error) {
	equipement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrEquipementNotFound
	}

	if req.Nom != nil {
		equipement.Nom = *req.Nom
	}

	if req.Code != nil {
		existing, _ := s.repo.GetByCode(ctx, *req.Code)
		if existing != nil && existing.IDEquipement != id {
			return nil, ErrCodeExists
		}
		equipement.Code = *req.Code
	}

	if req.Description != nil {
		equipement.Description = *req.Description
	}

	if req.Statut != nil {
		equipement.Statut = domain.EquipementStatus(*req.Statut)
	}

	if req.Localisation != nil {
		equipement.Localisation = *req.Localisation
	}

	if req.Criticite != nil {
		equipement.Criticite = *req.Criticite
	}

	if err := s.repo.Update(ctx, equipement); err != nil {
		return nil, fmt.Errorf("failed to update equipement: %w", err)
	}

	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload equipement: %w", err)
	}

	resp := updated.ToResponse()
	return &resp, nil
}

// DeleteEquipement removes an equipement by its UUID.
func (s *EquipementService) DeleteEquipement(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrEquipementNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete equipement: %w", err)
	}

	return nil
}
