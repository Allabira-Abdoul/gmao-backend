package service

import (
	"context"
	"time"

	"backend-gmao/apps/audit-service/internal/core/domain"
	"backend-gmao/apps/audit-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

// AuditService implements primary.AuditService.
type AuditService struct {
	auditRepo secondary.AuditLogRepository
}

// NewAuditService initializes a new AuditService instance.
func NewAuditService(auditRepo secondary.AuditLogRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

func (s *AuditService) WriteLog(ctx context.Context, req domain.CreateAuditLogRequest) (*domain.AuditLogResponse, error) {
	var userUUID *uuid.UUID
	if req.UserID != nil && *req.UserID != "" {
		parsed, err := uuid.Parse(*req.UserID)
		if err == nil {
			userUUID = &parsed
		}
	}

	auditLog := &domain.AuditLog{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Action:      req.Action,
		Details:     req.Details,
		UserID:      userUUID,
		PerformedAt: time.Now(),
	}

	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		return nil, err
	}

	resp := auditLog.ToResponse()
	return &resp, nil
}

func (s *AuditService) GetAllLogs(ctx context.Context) ([]domain.AuditLogResponse, error) {
	logs, err := s.auditRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.AuditLogResponse, len(logs))
	for i, l := range logs {
		responses[i] = l.ToResponse()
	}
	return responses, nil
}
