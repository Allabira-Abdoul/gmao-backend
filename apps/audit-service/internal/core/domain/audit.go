package domain

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents a secure log entry of any action taken across microservices.
type AuditLog struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ServiceName string     `gorm:"column:service_name;not null" json:"service_name"`
	Action      string     `gorm:"column:action;not null" json:"action"`
	Details     string     `gorm:"column:details" json:"details"`
	UserID      *uuid.UUID `gorm:"column:user_id;type:uuid" json:"user_id"`
	PerformedAt time.Time  `gorm:"column:performed_at;not null;default:now()" json:"performed_at"`
}

// TableName overrides GORM's default table name.
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogResponse represents the API DTO for AuditLog.
type AuditLogResponse struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Action      string     `json:"action"`
	Details     string     `json:"details"`
	UserID      *uuid.UUID `json:"user_id"`
	PerformedAt time.Time  `json:"performed_at"`
}

// ToResponse converts an AuditLog to AuditLogResponse DTO.
func (l *AuditLog) ToResponse() AuditLogResponse {
	return AuditLogResponse{
		ID:          l.ID,
		ServiceName: l.ServiceName,
		Action:      l.Action,
		Details:     l.Details,
		UserID:      l.UserID,
		PerformedAt: l.PerformedAt,
	}
}

// CreateAuditLogRequest is the DTO used to submit a new audit log.
type CreateAuditLogRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Action      string  `json:"action" binding:"required"`
	Details     string  `json:"details"`
	UserID      *string `json:"user_id,omitempty" binding:"omitempty,uuid"`
}
