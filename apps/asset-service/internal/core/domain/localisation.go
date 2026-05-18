package domain

import (
	"time"

	"github.com/google/uuid"
)

// Localisation represents a site, building, or room where equipment is placed.
type Localisation struct {
	IDLocalisation   uuid.UUID  `gorm:"column:id_localisation;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_localisation"`
	Nom              string     `gorm:"column:nom;not null" json:"nom"`
	Description      string     `gorm:"column:description" json:"description"`
	ParentID         *uuid.UUID `gorm:"column:parent_id;type:uuid" json:"parent_id"`
	LocalisationParent *Localisation `gorm:"foreignKey:ParentID;references:IDLocalisation;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent,omitempty"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (Localisation) TableName() string {
	return "localisations"
}

// LocalisationResponse is the DTO returned by API endpoints.
type LocalisationResponse struct {
	IDLocalisation uuid.UUID             `json:"id_localisation"`
	Nom            string                `json:"nom"`
	Description    string                `json:"description"`
	ParentID       *uuid.UUID            `json:"parent_id,omitempty"`
	Parent         *LocalisationResponse `json:"parent,omitempty"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

// ToResponse converts a Localisation to a LocalisationResponse.
func (l *Localisation) ToResponse() LocalisationResponse {
	resp := LocalisationResponse{
		IDLocalisation: l.IDLocalisation,
		Nom:            l.Nom,
		Description:    l.Description,
		ParentID:       l.ParentID,
		CreatedAt:      l.CreatedAt,
		UpdatedAt:      l.UpdatedAt,
	}

	if l.LocalisationParent != nil && l.LocalisationParent.IDLocalisation != uuid.Nil {
		parentResp := l.LocalisationParent.ToResponse()
		resp.Parent = &parentResp
	}

	return resp
}

// CreateLocalisationRequest is the DTO for creating a new location.
type CreateLocalisationRequest struct {
	Nom         string  `json:"nom" binding:"required,min=2,max=255"`
	Description string  `json:"description" binding:"omitempty"`
	ParentID    *string `json:"parent_id" binding:"omitempty,uuid"`
}

// UpdateLocalisationRequest is the DTO for updating an existing location.
type UpdateLocalisationRequest struct {
	Nom         *string `json:"nom,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty"`
	ParentID    *string `json:"parent_id,omitempty" binding:"omitempty,uuid"`
}
