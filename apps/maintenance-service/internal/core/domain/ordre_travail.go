package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrdreTravailStatut represents the status of a work order.
type OrdreTravailStatut string

const (
	StatutPlanifie OrdreTravailStatut = "PLANIFIE"
	StatutEnCours  OrdreTravailStatut = "EN_COURS"
	StatutTermine  OrdreTravailStatut = "TERMINE"
	StatutAnnule   OrdreTravailStatut = "ANNULE"
)

// TypeMaintenance represents the type of maintenance operation.
type TypeMaintenance string

const (
	TypePreventive TypeMaintenance = "PREVENTIVE"
	TypeCorrective TypeMaintenance = "CORRECTIVE"
)

// Priorite represents the priority level of a work order.
type Priorite string

const (
	PrioriteBasse   Priorite = "BASSE"
	PrioriteMoyenne Priorite = "MOYENNE"
	PrioriteHaute   Priorite = "HAUTE"
	PrioriteUrgente Priorite = "URGENTE"
)

// OrdreTravail represents a work order (Ordre de Travail) in the GMAO system.
type OrdreTravail struct {
	IDOrdreTravail       uuid.UUID          `gorm:"column:id_ordre_travail;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_ordre_travail"`
	Reference            string             `gorm:"column:reference;uniqueIndex;not null" json:"reference"`
	Titre                string             `gorm:"column:titre;not null" json:"titre"`
	Description          string             `gorm:"column:description" json:"description"`
	Statut               OrdreTravailStatut `gorm:"column:statut;type:varchar(50);default:'PLANIFIE'" json:"statut"`
	TypeMaintenance      TypeMaintenance    `gorm:"column:type_maintenance;type:varchar(20);not null" json:"type_maintenance"`
	Priorite             Priorite           `gorm:"column:priorite;type:varchar(20);default:'MOYENNE'" json:"priorite"`
	IDEquipement         uuid.UUID          `gorm:"column:id_equipement;type:uuid;not null" json:"id_equipement"`
	IDUtilisateurCree    uuid.UUID          `gorm:"column:id_utilisateur_cree;type:uuid;not null" json:"id_utilisateur_cree"`
	IDUtilisateurAssigne *uuid.UUID         `gorm:"column:id_utilisateur_assigne;type:uuid" json:"id_utilisateur_assigne"`
	IDEquipeAssigne      *uuid.UUID         `gorm:"column:id_equipe_assigne;type:uuid" json:"id_equipe_assigne"`
	IDCurrentState       *uuid.UUID         `gorm:"column:id_current_state;type:uuid" json:"id_current_state"`
	DateDebutPrevue      *time.Time         `gorm:"column:date_debut_prevue" json:"date_debut_prevue"`
	DateFinPrevue        *time.Time         `gorm:"column:date_fin_prevue" json:"date_fin_prevue"`
	DateDebutReelle      *time.Time         `gorm:"column:date_debut_reelle" json:"date_debut_reelle"`
	DateFinReelle        *time.Time         `gorm:"column:date_fin_reelle" json:"date_fin_reelle"`
	TempsPasse           int                `gorm:"column:temps_passe;default:0" json:"temps_passe"` // in minutes
	Commentaire          string             `gorm:"column:commentaire" json:"commentaire"`
	CreatedAt            time.Time          `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time          `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (OrdreTravail) TableName() string {
	return "ordres_travail"
}

// OrdreTravailResponse is the DTO returned by API endpoints.
type OrdreTravailResponse struct {
	IDOrdreTravail       uuid.UUID          `json:"id_ordre_travail"`
	Reference            string             `json:"reference"`
	Titre                string             `json:"titre"`
	Description          string             `json:"description"`
	Statut               OrdreTravailStatut `json:"statut"`
	TypeMaintenance      TypeMaintenance    `json:"type_maintenance"`
	Priorite             Priorite           `json:"priorite"`
	IDEquipement         uuid.UUID          `json:"id_equipement"`
	IDUtilisateurCree    uuid.UUID          `json:"id_utilisateur_cree"`
	IDUtilisateurAssigne *uuid.UUID         `json:"id_utilisateur_assigne"`
	IDEquipeAssigne      *uuid.UUID         `json:"id_equipe_assigne,omitempty"`
	IDCurrentState       *uuid.UUID         `json:"id_current_state,omitempty"`
	DateDebutPrevue      *time.Time         `json:"date_debut_prevue"`
	DateFinPrevue        *time.Time         `json:"date_fin_prevue"`
	DateDebutReelle      *time.Time         `json:"date_debut_reelle"`
	DateFinReelle        *time.Time         `json:"date_fin_reelle"`
	TempsPasse           int                `json:"temps_passe"`
	Commentaire          string             `json:"commentaire"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

// ToResponse converts an OrdreTravail to an OrdreTravailResponse.
func (o *OrdreTravail) ToResponse() OrdreTravailResponse {
	return OrdreTravailResponse{
		IDOrdreTravail:       o.IDOrdreTravail,
		Reference:            o.Reference,
		Titre:                o.Titre,
		Description:          o.Description,
		Statut:               o.Statut,
		TypeMaintenance:      o.TypeMaintenance,
		Priorite:             o.Priorite,
		IDEquipement:         o.IDEquipement,
		IDUtilisateurCree:    o.IDUtilisateurCree,
		IDUtilisateurAssigne: o.IDUtilisateurAssigne,
		IDEquipeAssigne:      o.IDEquipeAssigne,
		IDCurrentState:       o.IDCurrentState,
		DateDebutPrevue:      o.DateDebutPrevue,
		DateFinPrevue:        o.DateFinPrevue,
		DateDebutReelle:      o.DateDebutReelle,
		DateFinReelle:        o.DateFinReelle,
		TempsPasse:           o.TempsPasse,
		Commentaire:          o.Commentaire,
		CreatedAt:            o.CreatedAt,
		UpdatedAt:            o.UpdatedAt,
	}
}

// CreateOrdreTravailRequest is the DTO for creating a new work order.
type CreateOrdreTravailRequest struct {
	Titre                string  `json:"titre" binding:"required,min=2,max=255"`
	Description          string  `json:"description" binding:"omitempty"`
	TypeMaintenance      string  `json:"type_maintenance" binding:"required,oneof=PREVENTIVE CORRECTIVE"`
	Priorite             string  `json:"priorite" binding:"omitempty,oneof=BASSE MOYENNE HAUTE URGENTE"`
	IDEquipement         string  `json:"id_equipement" binding:"required,uuid"`
	IDUtilisateurAssigne *string `json:"id_utilisateur_assigne" binding:"omitempty,uuid"`
	IDEquipeAssigne      *string `json:"id_equipe_assigne" binding:"omitempty,uuid"`
	DateDebutPrevue      *string `json:"date_debut_prevue" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DateFinPrevue        *string `json:"date_fin_prevue" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Commentaire          string  `json:"commentaire" binding:"omitempty"`
}

// UpdateOrdreTravailRequest is the DTO for updating an existing work order.
type UpdateOrdreTravailRequest struct {
	Titre                *string `json:"titre,omitempty" binding:"omitempty,min=2,max=255"`
	Description          *string `json:"description,omitempty"`
	Statut               *string `json:"statut,omitempty" binding:"omitempty,oneof=PLANIFIE EN_COURS TERMINE ANNULE"`
	Priorite             *string `json:"priorite,omitempty" binding:"omitempty,oneof=BASSE MOYENNE HAUTE URGENTE"`
	IDUtilisateurAssigne *string `json:"id_utilisateur_assigne,omitempty" binding:"omitempty,uuid"`
	IDEquipeAssigne      *string `json:"id_equipe_assigne,omitempty" binding:"omitempty,uuid"`
	IDCurrentState       *string `json:"id_current_state,omitempty" binding:"omitempty,uuid"`
	DateDebutPrevue      *string `json:"date_debut_prevue,omitempty" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DateFinPrevue        *string `json:"date_fin_prevue,omitempty" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DateDebutReelle      *string `json:"date_debut_reelle,omitempty" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DateFinReelle        *string `json:"date_fin_reelle,omitempty" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Commentaire          *string `json:"commentaire,omitempty"`
}
