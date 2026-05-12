package domain

import (
	"time"

	"github.com/google/uuid"
)

// EquipementStatus represents the status of an equipment.
type EquipementStatus string

const (
	StatusEnService    EquipementStatus = "EN_SERVICE"
	StatusEnMaintenance EquipementStatus = "EN_MAINTENANCE"
	StatusEnPanne      EquipementStatus = "EN_PANNE"
	StatusReforme      EquipementStatus = "REFORME"
)

// Equipement represents a piece of machinery or tool in the GMAO system.
type Equipement struct {
	IDEquipement    uuid.UUID        `gorm:"column:id_equipement;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_equipement"`
	Nom             string           `gorm:"column:nom;not null" json:"nom"`
	Code            string           `gorm:"column:code;uniqueIndex;not null" json:"code"`
	Description     string           `gorm:"column:description" json:"description"`
	Statut          EquipementStatus `gorm:"column:statut;type:varchar(20);default:'EN_SERVICE'" json:"statut"`
	DateAcquisition time.Time        `gorm:"column:date_acquisition" json:"date_acquisition"`
	Localisation    string           `gorm:"column:localisation" json:"localisation"`
	Criticite       string           `gorm:"column:criticite;type:varchar(20);default:'BASSE'" json:"criticite"` // BASSE, MOYENNE, HAUTE
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
	Statut          EquipementStatus `json:"statut"`
	DateAcquisition time.Time        `json:"date_acquisition"`
	Localisation    string           `json:"localisation"`
	Criticite       string           `json:"criticite"`
	CreatedAt       time.Time        `json:"created_at"`
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
		Localisation:    e.Localisation,
		Criticite:       e.Criticite,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// CreateEquipementRequest is the DTO for creating a new equipment.
type CreateEquipementRequest struct {
	Nom             string `json:"nom" binding:"required,min=2,max=255"`
	Code            string `json:"code" binding:"required,min=2,max=50"`
	Description     string `json:"description" binding:"omitempty"`
	DateAcquisition string `json:"date_acquisition" binding:"required,datetime=2006-01-02"`
	Localisation    string `json:"localisation" binding:"omitempty,max=255"`
	Criticite       string `json:"criticite" binding:"omitempty,oneof=BASSE MOYENNE HAUTE"`
}

// UpdateEquipementRequest is the DTO for updating an existing equipment.
type UpdateEquipementRequest struct {
	Nom          *string `json:"nom,omitempty" binding:"omitempty,min=2,max=255"`
	Code         *string `json:"code,omitempty" binding:"omitempty,min=2,max=50"`
	Description  *string `json:"description,omitempty" binding:"omitempty"`
	Statut       *string `json:"statut,omitempty" binding:"omitempty,oneof=EN_SERVICE EN_MAINTENANCE EN_PANNE REFORME"`
	Localisation *string `json:"localisation,omitempty" binding:"omitempty,max=255"`
	Criticite    *string `json:"criticite,omitempty" binding:"omitempty,oneof=BASSE MOYENNE HAUTE"`
}
