package domain

import (
	"time"

	"github.com/google/uuid"
)

// MetricThreshold defines acceptable operational ranges for an asset or component.
type MetricThreshold struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AssetID     uuid.UUID  `gorm:"column:asset_id;type:uuid;not null" json:"asset_id"`
	ComponentID *uuid.UUID `gorm:"column:component_id;type:uuid" json:"component_id"`
	MetricName  string     `gorm:"column:metric_name;not null" json:"metric_name"`
	MinValue    *float64   `gorm:"column:min_value" json:"min_value"`
	MaxValue    *float64   `gorm:"column:max_value" json:"max_value"`
	Unit        string     `gorm:"column:unit;not null" json:"unit"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (MetricThreshold) TableName() string {
	return "metric_thresholds"
}

// MetricThresholdResponse is the DTO for MetricThreshold.
type MetricThresholdResponse struct {
	ID          uuid.UUID  `json:"id"`
	AssetID     uuid.UUID  `json:"asset_id"`
	ComponentID *uuid.UUID `json:"component_id"`
	MetricName  string     `json:"metric_name"`
	MinValue    *float64   `json:"min_value"`
	MaxValue    *float64   `json:"max_value"`
	Unit        string     `json:"unit"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (t *MetricThreshold) ToResponse() MetricThresholdResponse {
	return MetricThresholdResponse{
		ID:          t.ID,
		AssetID:     t.AssetID,
		ComponentID: t.ComponentID,
		MetricName:  t.MetricName,
		MinValue:    t.MinValue,
		MaxValue:    t.MaxValue,
		Unit:        t.Unit,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
