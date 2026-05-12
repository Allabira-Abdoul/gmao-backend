package domain

import (
	"time"

	"github.com/google/uuid"
)

// InterventionStatut represents the status of an intervention.
type InterventionStatut string

const (
	InterventionEnCours  InterventionStatut = "EN_COURS"
	InterventionTerminee InterventionStatut = "TERMINEE"
	InterventionAnnulee  InterventionStatut = "ANNULEE"
)

// Intervention represents a maintenance intervention linked to a work order.
// An OrdreTravail can have multiple Interventions (e.g., multiple visits to fix an issue).
type Intervention struct {
	IDIntervention    uuid.UUID          `gorm:"column:id_intervention;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_intervention"`
	IDOrdreTravail    uuid.UUID          `gorm:"column:id_ordre_travail;type:uuid;not null;index" json:"id_ordre_travail"`
	IDTechnicien      uuid.UUID          `gorm:"column:id_technicien;type:uuid;not null" json:"id_technicien"`
	Statut            InterventionStatut `gorm:"column:statut;type:varchar(20);default:'EN_COURS'" json:"statut"`
	DateDebut         time.Time          `gorm:"column:date_debut;not null" json:"date_debut"`
	DateFin           *time.Time         `gorm:"column:date_fin" json:"date_fin"`
	DureeMinutes      *int               `gorm:"column:duree_minutes" json:"duree_minutes"`
	RapportIntervention string           `gorm:"column:rapport_intervention;type:text" json:"rapport_intervention"`
	ActionsEffectuees string             `gorm:"column:actions_effectuees;type:text" json:"actions_effectuees"`
	CreatedAt         time.Time          `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time          `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (Intervention) TableName() string {
	return "interventions"
}

// InterventionResponse is the DTO returned by API endpoints.
type InterventionResponse struct {
	IDIntervention      uuid.UUID          `json:"id_intervention"`
	IDOrdreTravail      uuid.UUID          `json:"id_ordre_travail"`
	IDTechnicien        uuid.UUID          `json:"id_technicien"`
	Statut              InterventionStatut `json:"statut"`
	DateDebut           time.Time          `json:"date_debut"`
	DateFin             *time.Time         `json:"date_fin"`
	DureeMinutes        *int               `json:"duree_minutes"`
	RapportIntervention string             `json:"rapport_intervention"`
	ActionsEffectuees   string             `json:"actions_effectuees"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
}

// ToResponse converts an Intervention to an InterventionResponse.
func (i *Intervention) ToResponse() InterventionResponse {
	return InterventionResponse{
		IDIntervention:      i.IDIntervention,
		IDOrdreTravail:      i.IDOrdreTravail,
		IDTechnicien:        i.IDTechnicien,
		Statut:              i.Statut,
		DateDebut:           i.DateDebut,
		DateFin:             i.DateFin,
		DureeMinutes:        i.DureeMinutes,
		RapportIntervention: i.RapportIntervention,
		ActionsEffectuees:   i.ActionsEffectuees,
		CreatedAt:           i.CreatedAt,
		UpdatedAt:           i.UpdatedAt,
	}
}

// CreateInterventionRequest is the DTO for creating a new intervention.
type CreateInterventionRequest struct {
	IDOrdreTravail string `json:"id_ordre_travail" binding:"required,uuid"`
	IDTechnicien   string `json:"id_technicien" binding:"required,uuid"`
	DateDebut      string `json:"date_debut" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Commentaire    string `json:"commentaire" binding:"omitempty"`
}

// UpdateInterventionRequest is the DTO for updating an existing intervention.
type UpdateInterventionRequest struct {
	Statut              *string `json:"statut,omitempty" binding:"omitempty,oneof=EN_COURS TERMINEE ANNULEE"`
	DateFin             *string `json:"date_fin,omitempty" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DureeMinutes        *int    `json:"duree_minutes,omitempty" binding:"omitempty,min=1"`
	RapportIntervention *string `json:"rapport_intervention,omitempty"`
	ActionsEffectuees   *string `json:"actions_effectuees,omitempty"`
}
