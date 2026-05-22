package domain

import (
	"time"

	"github.com/google/uuid"
)

// AssetComponent represents a sub-part of an asset.
type AssetComponent struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AssetID     uuid.UUID `gorm:"column:asset_id;type:uuid;not null" json:"asset_id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (AssetComponent) TableName() string {
	return "asset_components"
}

// AssetComponentResponse is the DTO for AssetComponent.
type AssetComponentResponse struct {
	ID          uuid.UUID `json:"id"`
	AssetID     uuid.UUID `json:"asset_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *AssetComponent) ToResponse() AssetComponentResponse {
	return AssetComponentResponse{
		ID:          c.ID,
		AssetID:     c.AssetID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
