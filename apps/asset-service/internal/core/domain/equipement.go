package domain

import (
	"time"

	"github.com/google/uuid"
)

// EquipementStatus represents the status of an equipment.
type EquipementStatus string

const (
	StatusEnService     EquipementStatus = "EN_SERVICE"
	StatusEnMaintenance EquipementStatus = "EN_MAINTENANCE"
	StatusEnPanne       EquipementStatus = "EN_PANNE"
	StatusReforme       EquipementStatus = "REFORME"
)

// Equipement represents a piece of machinery or tool in the GMAO system.
type Equipement struct {
	IDEquipement    uuid.UUID        `gorm:"column:id_equipement;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_equipement"`
	Nom             string           `gorm:"column:nom;not null" json:"nom"`
	Code            string           `gorm:"column:code;uniqueIndex;not null" json:"code"`
	Description     string           `gorm:"column:description" json:"description"`
	Statut          EquipementStatus `gorm:"column:statut;type:varchar(20);default:'EN_SERVICE'" json:"statut"`
	DateAcquisition time.Time        `gorm:"column:date_acquisition" json:"date_acquisition"`
	IDLocalisation  *uuid.UUID       `gorm:"column:id_localisation;type:uuid" json:"id_localisation"`
	Localisation    *Localisation    `gorm:"foreignKey:IDLocalisation;references:IDLocalisation;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"localisation,omitempty"`
	Criticite       string           `gorm:"column:criticite;type:varchar(20);default:'BASSE'" json:"criticite"` // BASSE, MOYENNE, HAUTE
	TempsArret      int              `gorm:"column:temps_arret;default:0" json:"temps_arret"` // in minutes

	CreatedAt       time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (Equipement) TableName() string {
	return "equipements"
}

// EquipementResponse is the DTO returned by API endpoints.
type EquipementResponse struct {
	IDEquipement    uuid.UUID        `json:"id_equipement"`
	Nom             string           `json:"nom"`
	Code            string           `json:"code"`
	Description     string           `json:"description"`
	Statut          EquipementStatus      `json:"statut"`
	DateAcquisition time.Time             `json:"date_acquisition"`
	IDLocalisation  *uuid.UUID            `json:"id_localisation,omitempty"`
	Localisation    *LocalisationResponse `json:"localisation,omitempty"`
	Criticite       string                `json:"criticite"`
	TempsArret      int                   `json:"temps_arret"`
	CreatedAt       time.Time             `json:"created_at"`

	UpdatedAt       time.Time        `json:"updated_at"`
}

// ToResponse converts an Equipement to an EquipementResponse.
func (e *Equipement) ToResponse() EquipementResponse {
	return EquipementResponse{
		IDEquipement:    e.IDEquipement,
		Nom:             e.Nom,
		Code:            e.Code,
		Description:     e.Description,
		Statut:          e.Statut,
		DateAcquisition: e.DateAcquisition,
		IDLocalisation:  e.IDLocalisation,
		Criticite:       e.Criticite,
		TempsArret:      e.TempsArret,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}

	if e.Localisation != nil && e.Localisation.IDLocalisation != uuid.Nil {
		locResp := e.Localisation.ToResponse()
		resp.Localisation = &locResp
	}

	return resp
}

// CreateEquipementRequest is the DTO for creating a new equipment.
type CreateEquipementRequest struct {
	Nom             string `json:"nom" binding:"required,min=2,max=255"`
	Code            string `json:"code" binding:"required,min=2,max=50"`
	Description     string `json:"description" binding:"omitempty"`
	DateAcquisition string `json:"date_acquisition" binding:"required,datetime=2006-01-02"`
	IDLocalisation  string `json:"id_localisation" binding:"omitempty,uuid"`
	Criticite       string `json:"criticite" binding:"omitempty,oneof=BASSE MOYENNE HAUTE"`
}

// UpdateEquipementRequest is the DTO for updating an existing equipment.
type UpdateEquipementRequest struct {
	Nom          *string `json:"nom,omitempty" binding:"omitempty,min=2,max=255"`
	Code         *string `json:"code,omitempty" binding:"omitempty,min=2,max=50"`
	Description    *string `json:"description,omitempty" binding:"omitempty"`
	Statut         *string `json:"statut,omitempty" binding:"omitempty,oneof=EN_SERVICE EN_MAINTENANCE EN_PANNE REFORME"`
	IDLocalisation *string `json:"id_localisation,omitempty" binding:"omitempty,uuid"`
	Criticite      *string `json:"criticite,omitempty" binding:"omitempty,oneof=BASSE MOYENNE HAUTE"`
	TempsArret     *int    `json:"temps_arret,omitempty" binding:"omitempty,min=0"`
}
