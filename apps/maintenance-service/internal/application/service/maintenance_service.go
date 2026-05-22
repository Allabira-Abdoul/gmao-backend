package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"backend-gmao/apps/maintenance-service/internal/core/ports/secondary"
	"backend-gmao/pkg/audit"
	"backend-gmao/pkg/middleware"

	"github.com/google/uuid"
)

var (
	ErrWorkOrderNotFound       = errors.New("work order not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// MaintenanceService implements primary.MaintenanceService.
type MaintenanceService struct {
	maintenanceRepo secondary.MaintenanceRepository
	analyticsClient secondary.AnalyticsClient
	auditClient     audit.Client
}

// NewMaintenanceService initializes a new MaintenanceService instance.
func NewMaintenanceService(
	maintenanceRepo secondary.MaintenanceRepository,
	analyticsClient secondary.AnalyticsClient,
	auditClient audit.Client,
) *MaintenanceService {
	return &MaintenanceService{
		maintenanceRepo: maintenanceRepo,
		analyticsClient: analyticsClient,
		auditClient:     auditClient,
	}
}

func (s *MaintenanceService) fireAudit(ctx context.Context, action, details string) {
	userID, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	var uidPtr *string
	if ok && userID != "" {
		uidPtr = &userID
	}

	go func() {
		bgCtx := context.Background()
		_ = s.auditClient.LogEvent(bgCtx, audit.AuditEvent{
			ServiceName: "maintenance-service",
			Action:      action,
			Details:     details,
			UserID:      uidPtr,
		})
	}()
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
		ID:                  uuid.New(),
		Title:               req.Title,
		Description:         req.Description,
		AssetID:             assetID,
		Priority:            req.Priority,
		Status:              "PENDING",
		MaintenanceCategory: req.MaintenanceCategory,
		MaintenanceType:     req.MaintenanceType,
		IsMetricMeasurement: req.IsMetricMeasurement,
		AssignedTo:          assignedTo,
	}

	if err := s.maintenanceRepo.CreateWorkOrder(ctx, wo); err != nil {
		return nil, err
	}

	s.fireAudit(ctx, "CREATE_WORK_ORDER", fmt.Sprintf("Created work order %s for asset %s", wo.ID, wo.AssetID))

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
	if req.MaintenanceCategory != nil {
		wo.MaintenanceCategory = *req.MaintenanceCategory
	}
	if req.MaintenanceType != nil {
		wo.MaintenanceType = *req.MaintenanceType
	}
	if req.IsMetricMeasurement != nil {
		wo.IsMetricMeasurement = *req.IsMetricMeasurement
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

	s.fireAudit(ctx, "UPDATE_WORK_ORDER", fmt.Sprintf("Updated work order %s", wo.ID))

	resp := wo.ToResponse()
	return &resp, nil
}

func (s *MaintenanceService) DeleteWorkOrder(ctx context.Context, id uuid.UUID) error {
	_, err := s.maintenanceRepo.FindWorkOrderByID(ctx, id)
	if err != nil {
		return ErrWorkOrderNotFound
	}
	
	if err := s.maintenanceRepo.DeleteWorkOrder(ctx, id); err != nil {
		return err
	}
	
	s.fireAudit(ctx, "DELETE_WORK_ORDER", fmt.Sprintf("Deleted work order %s", id))
	return nil
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
	wo, err := s.maintenanceRepo.FindWorkOrderByID(ctx, workOrderID)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	performedBy, err := uuid.Parse(req.PerformedBy)
	if err != nil {
		return nil, errors.New("invalid performed by user ID format")
	}

	intervention := &domain.Intervention{
		ID:                  uuid.New(),
		WorkOrderID:         workOrderID,
		Description:         req.Description,
		MaintenanceCategory: req.MaintenanceCategory,
		MaintenanceType:     req.MaintenanceType,
		IsMetricMeasurement: req.IsMetricMeasurement,
		DurationMinutes:     req.DurationMinutes,
		PerformedBy:         performedBy,
	}

	if len(req.Measurements) > 0 {
		meas := make([]domain.MetricMeasurement, len(req.Measurements))
		for i, mReq := range req.Measurements {
			var compID *uuid.UUID
			if mReq.ComponentID != nil && *mReq.ComponentID != "" {
				parsedComp, err := uuid.Parse(*mReq.ComponentID)
				if err == nil {
					compID = &parsedComp
				}
			}
			meas[i] = domain.MetricMeasurement{
				ID:                  uuid.New(),
				InterventionID:      intervention.ID,
				ComponentID:         compID,
				MetricName:          mReq.MetricName,
				Value:               mReq.Value,
				Unit:                mReq.Unit,
				IsThresholdBreached: mReq.IsThresholdBreached,
			}
		}
		intervention.Measurements = meas
	}

	if err := s.maintenanceRepo.CreateIntervention(ctx, intervention); err != nil {
		return nil, err
	}

	// Trigger analytics event for all interventions
	event := secondary.MaintenanceEvent{
		AssetID:             wo.AssetID,
		MaintenanceCategory: intervention.MaintenanceCategory,
		DurationMinutes:     float64(intervention.DurationMinutes),
	}
	
	// We run this asynchronously so it doesn't block the API response
	go func() {
		// Create a background context since the request context might be cancelled
		bgCtx := context.Background()
		_ = s.analyticsClient.PublishMaintenanceEvent(bgCtx, event)
	}()

	s.fireAudit(ctx, "RECORD_INTERVENTION", fmt.Sprintf("Recorded intervention %s for work order %s", intervention.ID, intervention.WorkOrderID))

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
