package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrdreTravail represents a work order in the GMAO system.
type OrdreTravail struct {
	ID                  uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title               string         `gorm:"column:title;not null" json:"title"`
	Description         string         `gorm:"column:description" json:"description"`
	AssetID             uuid.UUID      `gorm:"column:asset_id;type:uuid;not null" json:"asset_id"`
	Priority            string         `gorm:"column:priority;not null;default:'MEDIUM'" json:"priority"` // LOW, MEDIUM, HIGH, CRITICAL
	Status              string         `gorm:"column:status;not null;default:'PENDING'" json:"status"`    // PENDING, IN_PROGRESS, COMPLETED, CANCELLED
	MaintenanceCategory string         `gorm:"column:maintenance_category" json:"maintenance_category"`   // CORRECTIVE, PREVENTIVE
	MaintenanceType     string         `gorm:"column:maintenance_type" json:"maintenance_type"`           // PALLIATIVE, CURATIVE, SYSTEMATIC, CONDITIONAL, PREDICTIVE
	IsMetricMeasurement bool           `gorm:"column:is_metric_measurement;not null;default:false" json:"is_metric_measurement"`
	AssignedTo          *uuid.UUID     `gorm:"column:assigned_to;type:uuid" json:"assigned_to"`
	Interventions       []Intervention `gorm:"foreignKey:WorkOrderID" json:"interventions,omitempty"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides GORM's default table name.
func (OrdreTravail) TableName() string {
	return "work_orders"
}

// Intervention represents an action taken on a Work Order.
type Intervention struct {
	ID                  uuid.UUID           `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WorkOrderID         uuid.UUID           `gorm:"column:work_order_id;type:uuid;not null" json:"work_order_id"`
	Description         string              `gorm:"column:description;not null" json:"description"`
	MaintenanceCategory string              `gorm:"column:maintenance_category" json:"maintenance_category"`
	MaintenanceType     string              `gorm:"column:maintenance_type" json:"maintenance_type"`
	IsMetricMeasurement bool                `gorm:"column:is_metric_measurement;not null;default:false" json:"is_metric_measurement"`
	DurationMinutes     int                 `gorm:"column:duration_minutes;not null" json:"duration_minutes"`
	PerformedBy         uuid.UUID           `gorm:"column:performed_by;type:uuid;not null" json:"performed_by"`
	Measurements        []MetricMeasurement `gorm:"foreignKey:InterventionID" json:"measurements,omitempty"`
	CreatedAt           time.Time           `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time           `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides GORM's default table name.
func (Intervention) TableName() string {
	return "interventions"
}

// OrdreTravailResponse represents the DTO returned by API endpoints.
type OrdreTravailResponse struct {
	ID                  uuid.UUID              `json:"id"`
	Title               string                 `json:"title"`
	Description         string                 `json:"description"`
	AssetID             uuid.UUID              `json:"asset_id"`
	Priority            string                 `json:"priority"`
	Status              string                 `json:"status"`
	MaintenanceCategory string                 `json:"maintenance_category"`
	MaintenanceType     string                 `json:"maintenance_type"`
	IsMetricMeasurement bool                   `json:"is_metric_measurement"`
	AssignedTo          *uuid.UUID             `json:"assigned_to"`
	Interventions       []InterventionResponse `json:"interventions,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// InterventionResponse represents the API DTO for Intervention.
type InterventionResponse struct {
	ID                  uuid.UUID                   `json:"id"`
	WorkOrderID         uuid.UUID                   `json:"work_order_id"`
	Description         string                      `json:"description"`
	MaintenanceCategory string                      `json:"maintenance_category"`
	MaintenanceType     string                      `json:"maintenance_type"`
	IsMetricMeasurement bool                        `json:"is_metric_measurement"`
	DurationMinutes     int                         `json:"duration_minutes"`
	PerformedBy         uuid.UUID                   `json:"performed_by"`
	Measurements        []MetricMeasurementResponse `json:"measurements,omitempty"`
	CreatedAt           time.Time                   `json:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at"`
}

// ToResponse converts an Intervention to InterventionResponse.
func (i *Intervention) ToResponse() InterventionResponse {
	measResp := make([]MetricMeasurementResponse, len(i.Measurements))
	for idx, m := range i.Measurements {
		measResp[idx] = m.ToResponse()
	}

	return InterventionResponse{
		ID:                  i.ID,
		WorkOrderID:         i.WorkOrderID,
		Description:         i.Description,
		MaintenanceCategory: i.MaintenanceCategory,
		MaintenanceType:     i.MaintenanceType,
		IsMetricMeasurement: i.IsMetricMeasurement,
		DurationMinutes:     i.DurationMinutes,
		PerformedBy:         i.PerformedBy,
		Measurements:        measResp,
		CreatedAt:           i.CreatedAt,
		UpdatedAt:           i.UpdatedAt,
	}
}

// ToResponse converts an OrdreTravail to OrdreTravailResponse.
func (o *OrdreTravail) ToResponse() OrdreTravailResponse {
	interventionsResp := make([]InterventionResponse, len(o.Interventions))
	for idx, i := range o.Interventions {
		interventionsResp[idx] = i.ToResponse()
	}

	return OrdreTravailResponse{
		ID:                  o.ID,
		Title:               o.Title,
		Description:         o.Description,
		AssetID:             o.AssetID,
		Priority:            o.Priority,
		Status:              o.Status,
		MaintenanceCategory: o.MaintenanceCategory,
		MaintenanceType:     o.MaintenanceType,
		IsMetricMeasurement: o.IsMetricMeasurement,
		AssignedTo:          o.AssignedTo,
		Interventions:       interventionsResp,
		CreatedAt:           o.CreatedAt,
		UpdatedAt:           o.UpdatedAt,
	}
}

// CreateOrdreTravailRequest is the DTO to create a new work order.
type CreateOrdreTravailRequest struct {
	Title               string  `json:"title" binding:"required,min=2,max=255"`
	Description         string  `json:"description"`
	AssetID             string  `json:"asset_id" binding:"required,uuid"`
	Priority            string  `json:"priority" binding:"required"`
	MaintenanceCategory string  `json:"maintenance_category,omitempty"`
	MaintenanceType     string  `json:"maintenance_type,omitempty"`
	IsMetricMeasurement bool    `json:"is_metric_measurement,omitempty"`
	AssignedTo          *string `json:"assigned_to,omitempty" binding:"omitempty,uuid"`
}

// UpdateOrdreTravailRequest is the DTO to update an existing work order.
type UpdateOrdreTravailRequest struct {
	Title               *string `json:"title,omitempty" binding:"omitempty,min=2,max=255"`
	Description         *string `json:"description,omitempty"`
	Status              *string `json:"status,omitempty" binding:"omitempty,oneof=PENDING IN_PROGRESS COMPLETED CANCELLED"`
	Priority            *string `json:"priority,omitempty" binding:"omitempty,oneof=LOW MEDIUM HIGH CRITICAL"`
	MaintenanceCategory *string `json:"maintenance_category,omitempty"`
	MaintenanceType     *string `json:"maintenance_type,omitempty"`
	IsMetricMeasurement *bool   `json:"is_metric_measurement,omitempty"`
	AssignedTo          *string `json:"assigned_to,omitempty" binding:"omitempty,uuid"`
}

type CreateMetricMeasurementRequest struct {
	ComponentID         *string `json:"component_id,omitempty" binding:"omitempty,uuid"`
	MetricName          string  `json:"metric_name" binding:"required"`
	Value               float64 `json:"value" binding:"required"`
	Unit                string  `json:"unit" binding:"required"`
	IsThresholdBreached bool    `json:"is_threshold_breached"`
}

// CreateInterventionRequest is the DTO to record a new intervention.
type CreateInterventionRequest struct {
	Description         string                           `json:"description" binding:"required,min=2"`
	MaintenanceCategory string                           `json:"maintenance_category,omitempty"`
	MaintenanceType     string                           `json:"maintenance_type,omitempty"`
	IsMetricMeasurement bool                             `json:"is_metric_measurement,omitempty"`
	DurationMinutes     int                              `json:"duration_minutes" binding:"required,gt=0"`
	PerformedBy         string                           `json:"performed_by" binding:"required,uuid"`
	Measurements        []CreateMetricMeasurementRequest `json:"measurements,omitempty" binding:"omitempty,dive"`
}
