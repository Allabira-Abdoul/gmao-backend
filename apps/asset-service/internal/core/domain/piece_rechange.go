package domain

import (
	"time"

	"github.com/google/uuid"
)

// PieceRechange represents a spare part in the GMAO system.
type PieceRechange struct {
	IDPiece         uuid.UUID `gorm:"column:id_piece;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_piece"`
	Nom             string    `gorm:"column:nom;not null" json:"nom"`
	Reference       string    `gorm:"column:reference;uniqueIndex;not null" json:"reference"`
	Description     string    `gorm:"column:description" json:"description"`
	QuantiteEnStock int       `gorm:"column:quantite_en_stock;default:0;not null" json:"quantite_en_stock"`
	SeuilAlerte     int       `gorm:"column:seuil_alerte;default:0;not null" json:"seuil_alerte"`
	CoutUnitaire    float64   `gorm:"column:cout_unitaire;type:numeric(10,2);default:0.00" json:"cout_unitaire"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (PieceRechange) TableName() string {
	return "pieces_rechange"
}

// PieceRechangeResponse is the DTO returned by API endpoints.
type PieceRechangeResponse struct {
	IDPiece         uuid.UUID `json:"id_piece"`
	Nom             string    `json:"nom"`
	Reference       string    `json:"reference"`
	Description     string    `json:"description"`
	QuantiteEnStock int       `json:"quantite_en_stock"`
	SeuilAlerte     int       `json:"seuil_alerte"`
	CoutUnitaire    float64   `json:"cout_unitaire"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ToResponse converts a PieceRechange to a PieceRechangeResponse.
func (p *PieceRechange) ToResponse() PieceRechangeResponse {
	return PieceRechangeResponse{
		IDPiece:         p.IDPiece,
		Nom:             p.Nom,
		Reference:       p.Reference,
		Description:     p.Description,
		QuantiteEnStock: p.QuantiteEnStock,
		SeuilAlerte:     p.SeuilAlerte,
		CoutUnitaire:    p.CoutUnitaire,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// CreatePieceRechangeRequest is the DTO for creating a new spare part.
type CreatePieceRechangeRequest struct {
	Nom             string  `json:"nom" binding:"required,min=2,max=255"`
	Reference       string  `json:"reference" binding:"required,min=2,max=100"`
	Description     string  `json:"description" binding:"omitempty"`
	QuantiteEnStock int     `json:"quantite_en_stock" binding:"omitempty,min=0"`
	SeuilAlerte     int     `json:"seuil_alerte" binding:"omitempty,min=0"`
	CoutUnitaire    float64 `json:"cout_unitaire" binding:"omitempty,min=0"`
}

// UpdatePieceRechangeRequest is the DTO for updating an existing spare part.
type UpdatePieceRechangeRequest struct {
	Nom             *string  `json:"nom,omitempty" binding:"omitempty,min=2,max=255"`
	Reference       *string  `json:"reference,omitempty" binding:"omitempty,min=2,max=100"`
	Description     *string  `json:"description,omitempty" binding:"omitempty"`
	QuantiteEnStock *int     `json:"quantite_en_stock,omitempty" binding:"omitempty,min=0"`
	SeuilAlerte     *int     `json:"seuil_alerte,omitempty" binding:"omitempty,min=0"`
	CoutUnitaire    *float64 `json:"cout_unitaire,omitempty" binding:"omitempty,min=0"`
}
