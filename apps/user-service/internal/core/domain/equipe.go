package domain

import (
	"time"

	"github.com/google/uuid"
)

// Equipe represents a team or group of users in the GMAO system.
type Equipe struct {
	IDEquipe    uuid.UUID `gorm:"column:id_equipe;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_equipe"`
	Nom         string    `gorm:"column:nom;uniqueIndex;not null" json:"nom"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (Equipe) TableName() string {
	return "equipes"
}

// EquipeResponse is the DTO returned by API endpoints.
type EquipeResponse struct {
	IDEquipe    uuid.UUID `json:"id_equipe"`
	Nom         string    `json:"nom"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts an Equipe to an EquipeResponse.
func (e *Equipe) ToResponse() EquipeResponse {
	return EquipeResponse{
		IDEquipe:    e.IDEquipe,
		Nom:         e.Nom,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// CreateEquipeRequest is the DTO for creating a new team.
type CreateEquipeRequest struct {
	Nom         string `json:"nom" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"omitempty"`
}

// UpdateEquipeRequest is the DTO for updating an existing team.
type UpdateEquipeRequest struct {
	Nom         *string `json:"nom,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty"`
}
