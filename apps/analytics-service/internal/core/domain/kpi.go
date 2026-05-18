package domain

import (
	"time"

	"github.com/google/uuid"
)

// KPIMetric represents a generic key performance indicator structure
type KPIMetric struct {
	IDMetric  uuid.UUID `gorm:"column:id_metric;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_metric"`
	Name      string    `gorm:"column:name;not null" json:"name"` // e.g., "MTBF", "MTTR", "COMPLIANCE_RATE"
	Value     float64   `gorm:"column:value;not null" json:"value"`
	Unit      string    `gorm:"column:unit" json:"unit"` // e.g., "hours", "%", "USD"
	Period    string    `gorm:"column:period" json:"period"` // e.g., "2023-10", "2023-Q4"
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (KPIMetric) TableName() string {
	return "kpi_metrics"
}

// EquipmentDowntime tracks the downtime aggregated for each equipment
type EquipmentDowntime struct {
	IDEquipement uuid.UUID `gorm:"column:id_equipement;type:uuid;primaryKey" json:"id_equipement"`
	TotalDowntime int       `gorm:"column:total_downtime;default:0" json:"total_downtime"` // in minutes
	IncidentCount int       `gorm:"column:incident_count;default:0" json:"incident_count"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (EquipmentDowntime) TableName() string {
	return "equipment_downtime_reports"
}

// CostAnalysis tracks costs related to maintenance and spare parts
type CostAnalysis struct {
	IDPeriod   string    `gorm:"column:id_period;primaryKey" json:"id_period"` // e.g., "2023-10"
	LaborCost  float64   `gorm:"column:labor_cost;default:0.0" json:"labor_cost"`
	PartsCost  float64   `gorm:"column:parts_cost;default:0.0" json:"parts_cost"`
	TotalCost  float64   `gorm:"column:total_cost;default:0.0" json:"total_cost"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (CostAnalysis) TableName() string {
	return "cost_analysis_reports"
}
