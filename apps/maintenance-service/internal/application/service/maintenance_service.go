package service

import (
	"context"
	"errors"
	"time"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/apps/maintenance-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

var (
	ErrWorkOrderNotFound = errors.New("work order not found")
)

// MaintenanceService implements primary.MaintenanceService.
type MaintenanceService struct {
	maintenanceRepo secondary.MaintenanceRepository
}

// NewMaintenanceService initializes a new MaintenanceService instance.
func NewMaintenanceService(maintenanceRepo secondary.MaintenanceRepository) *MaintenanceService {
	return &MaintenanceService{maintenanceRepo: maintenanceRepo}
}

func (s *MaintenanceService) CreateWorkOrder(ctx context.Context, req domain.CreateOrdreTravailRequest) (*domain.OrdreTravailResponse, error) {
	assetID, err := uuid.Parse(req.AssetID)
	if err != nil {
		return nil, errors.New("invalid asset ID format")
	}

	var assignedTo *uuid.UUID
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		parsed, err := uuid.Parse(*req.AssignedTo)
		if err != nil {
			return nil, errors.New("invalid assigned user ID format")
		}
		assignedTo = &parsed
	}

	wo := &domain.OrdreTravail{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		AssetID:     assetID,
		Priority:    req.Priority,
		Status:      "PENDING",
		AssignedTo:  assignedTo,
	}

	if err := s.maintenanceRepo.CreateWorkOrder(ctx, wo); err != nil {
		return nil, err
	}

	resp := wo.ToResponse()
	return &resp, nil
}

func (s *MaintenanceService) UpdateWorkOrder(ctx context.Context, id uuid.UUID, req domain.UpdateOrdreTravailRequest) (*domain.OrdreTravailResponse, error) {
	wo, err := s.maintenanceRepo.FindWorkOrderByID(ctx, id)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	if req.Title != nil {
		wo.Title = *req.Title
	}
	if req.Description != nil {
		wo.Description = *req.Description
	}
	if req.Status != nil {
		wo.Status = *req.Status
	}
	if req.Priority != nil {
		wo.Priority = *req.Priority
	}
	if req.AssignedTo != nil {
		if *req.AssignedTo == "" {
			wo.AssignedTo = nil
		} else {
			parsed, err := uuid.Parse(*req.AssignedTo)
			if err != nil {
				return nil, errors.New("invalid assigned user ID format")
			}
			wo.AssignedTo = &parsed
		}
	}

	wo.UpdatedAt = time.Now()

	if err := s.maintenanceRepo.UpdateWorkOrder(ctx, wo); err != nil {
		return nil, err
	}

	resp := wo.ToResponse()
	return &resp, nil
}

func (s *MaintenanceService) DeleteWorkOrder(ctx context.Context, id uuid.UUID) error {
	_, err := s.maintenanceRepo.FindWorkOrderByID(ctx, id)
	if err != nil {
		return ErrWorkOrderNotFound
	}
	return s.maintenanceRepo.DeleteWorkOrder(ctx, id)
}

func (s *MaintenanceService) GetWorkOrder(ctx context.Context, id uuid.UUID) (*domain.OrdreTravailResponse, error) {
	wo, err := s.maintenanceRepo.FindWorkOrderByID(ctx, id)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	interventions, _ := s.maintenanceRepo.FindInterventionsByWorkOrderID(ctx, id)
	wo.Interventions = interventions

	resp := wo.ToResponse()
	return &resp, nil
}

func (s *MaintenanceService) GetAllWorkOrders(ctx context.Context) ([]domain.OrdreTravailResponse, error) {
	workorders, err := s.maintenanceRepo.FindAllWorkOrders(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.OrdreTravailResponse, len(workorders))
	for i, wo := range workorders {
		interventions, _ := s.maintenanceRepo.FindInterventionsByWorkOrderID(ctx, wo.ID)
		wo.Interventions = interventions
		responses[i] = wo.ToResponse()
	}
	return responses, nil
}

func (s *MaintenanceService) RecordIntervention(ctx context.Context, workOrderID uuid.UUID, req domain.CreateInterventionRequest) (*domain.InterventionResponse, error) {
	_, err := s.maintenanceRepo.FindWorkOrderByID(ctx, workOrderID)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	performedBy, err := uuid.Parse(req.PerformedBy)
	if err != nil {
		return nil, errors.New("invalid performed by user ID format")
	}

	intervention := &domain.Intervention{
		ID:              uuid.New(),
		WorkOrderID:     workOrderID,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		PerformedBy:     performedBy,
	}

	if err := s.maintenanceRepo.CreateIntervention(ctx, intervention); err != nil {
		return nil, err
	}

	resp := intervention.ToResponse()
	return &resp, nil
}

func (s *MaintenanceService) GetInterventionsForWorkOrder(ctx context.Context, workOrderID uuid.UUID) ([]domain.InterventionResponse, error) {
	interventions, err := s.maintenanceRepo.FindInterventionsByWorkOrderID(ctx, workOrderID)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.InterventionResponse, len(interventions))
	for i, iv := range interventions {
		responses[i] = iv.ToResponse()
	}
	return responses, nil
}
