package domain

import (
	"time"

	"github.com/google/uuid"
)

// StatutCommande represents the status of a purchase order.
type StatutCommande string

const (
	StatutBrouillon StatutCommande = "BROUILLON"
	StatutEnAttente StatutCommande = "EN_ATTENTE_VALIDATION"
	StatutValidee   StatutCommande = "VALIDEE"
	StatutRefusee   StatutCommande = "REFUSEE"
	StatutRecue     StatutCommande = "RECUE"
)

// CommandeAchat represents a purchase order for spare parts.
type CommandeAchat struct {
	IDCommande      uuid.UUID      `gorm:"column:id_commande;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_commande"`
	Reference       string         `gorm:"column:reference;uniqueIndex;not null" json:"reference"`
	IDPieceRechange uuid.UUID      `gorm:"column:id_piece_rechange;type:uuid;not null" json:"id_piece_rechange"`
	Quantite        int            `gorm:"column:quantite;not null" json:"quantite"`
	PrixUnitaire    float64        `gorm:"column:prix_unitaire" json:"prix_unitaire"`
	Statut          StatutCommande `gorm:"column:statut;type:varchar(30);default:'BROUILLON'" json:"statut"`
	IDUtilisateur   uuid.UUID      `gorm:"column:id_utilisateur;type:uuid;not null" json:"id_utilisateur"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (CommandeAchat) TableName() string {
	return "commandes_achat"
}

// CommandeAchatResponse is the DTO returned by API endpoints.
type CommandeAchatResponse struct {
	IDCommande      uuid.UUID      `json:"id_commande"`
	Reference       string         `json:"reference"`
	IDPieceRechange uuid.UUID      `json:"id_piece_rechange"`
	Quantite        int            `json:"quantite"`
	PrixUnitaire    float64        `json:"prix_unitaire"`
	Statut          StatutCommande `json:"statut"`
	IDUtilisateur   uuid.UUID      `json:"id_utilisateur"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// ToResponse converts a CommandeAchat to a CommandeAchatResponse.
func (c *CommandeAchat) ToResponse() CommandeAchatResponse {
	return CommandeAchatResponse{
		IDCommande:      c.IDCommande,
		Reference:       c.Reference,
		IDPieceRechange: c.IDPieceRechange,
		Quantite:        c.Quantite,
		PrixUnitaire:    c.PrixUnitaire,
		Statut:          c.Statut,
		IDUtilisateur:   c.IDUtilisateur,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

// CreateCommandeAchatRequest is the DTO for creating a new purchase order.
type CreateCommandeAchatRequest struct {
	IDPieceRechange string  `json:"id_piece_rechange" binding:"required,uuid"`
	Quantite        int     `json:"quantite" binding:"required,min=1"`
	PrixUnitaire    float64 `json:"prix_unitaire" binding:"omitempty,min=0"`
}

// UpdateCommandeAchatRequest is the DTO for updating an existing purchase order.
type UpdateCommandeAchatRequest struct {
	Quantite     *int    `json:"quantite,omitempty" binding:"omitempty,min=1"`
	PrixUnitaire *float64 `json:"prix_unitaire,omitempty" binding:"omitempty,min=0"`
	Statut       *string `json:"statut,omitempty" binding:"omitempty,oneof=BROUILLON EN_ATTENTE_VALIDATION VALIDEE REFUSEE RECUE"`
}
