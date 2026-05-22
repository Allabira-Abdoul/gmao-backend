package domain

import (
	"time"

	"github.com/google/uuid"
)

// Asset represents a physical asset in the GMAO system.
type Asset struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	Code          string    `gorm:"column:code;uniqueIndex;not null" json:"code"`
	Status        string    `gorm:"column:status;not null;default:'OPERATIONAL'" json:"status"` // OPERATIONAL, DOWN, UNDER_REPAIR, SCRAPPED
	Category      string    `gorm:"column:category;not null" json:"category"`
	Location      string    `gorm:"column:location;not null" json:"location"`
	PurchaseDate  time.Time `gorm:"column:purchase_date" json:"purchase_date"`
	PurchaseValue float64           `gorm:"column:purchase_value" json:"purchase_value"`
	Components    []AssetComponent  `gorm:"foreignKey:AssetID" json:"components,omitempty"`
	Thresholds    []MetricThreshold `gorm:"foreignKey:AssetID" json:"thresholds,omitempty"`
	CreatedAt     time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides GORM's default table name.
func (Asset) TableName() string {
	return "assets"
}

// AssetResponse represents the API DTO for Asset.
type AssetResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Status        string    `json:"status"`
	Category      string    `json:"category"`
	Location      string    `json:"location"`
	PurchaseDate  time.Time                 `json:"purchase_date"`
	PurchaseValue float64                   `json:"purchase_value"`
	Components    []AssetComponentResponse  `json:"components,omitempty"`
	Thresholds    []MetricThresholdResponse `json:"thresholds,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

// ToResponse converts an Asset to AssetResponse DTO.
func (a *Asset) ToResponse() AssetResponse {
	comps := make([]AssetComponentResponse, len(a.Components))
	for i, c := range a.Components {
		comps[i] = c.ToResponse()
	}

	thresh := make([]MetricThresholdResponse, len(a.Thresholds))
	for i, t := range a.Thresholds {
		thresh[i] = t.ToResponse()
	}

	return AssetResponse{
		ID:            a.ID,
		Name:          a.Name,
		Code:          a.Code,
		Status:        a.Status,
		Category:      a.Category,
		Location:      a.Location,
		PurchaseDate:  a.PurchaseDate,
		PurchaseValue: a.PurchaseValue,
		Components:    comps,
		Thresholds:    thresh,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

// CreateAssetRequest is the DTO used to submit a new asset.
type CreateAssetRequest struct {
	Name          string    `json:"name" binding:"required,min=2,max=255"`
	Code          string    `json:"code" binding:"required,min=2,max=50"`
	Category      string    `json:"category" binding:"required"`
	Location      string    `json:"location" binding:"required"`
	PurchaseDate  time.Time `json:"purchase_date" binding:"required"`
	PurchaseValue float64   `json:"purchase_value" binding:"required"`
}

// UpdateAssetRequest is the DTO used to update an existing asset.
type UpdateAssetRequest struct {
	Name          *string  `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	Status        *string  `json:"status,omitempty" binding:"omitempty,oneof=OPERATIONAL DOWN UNDER_REPAIR SCRAPPED"`
	Category      *string  `json:"category,omitempty"`
	Location      *string  `json:"location,omitempty"`
	PurchaseValue *float64 `json:"purchase_value,omitempty"`
}
