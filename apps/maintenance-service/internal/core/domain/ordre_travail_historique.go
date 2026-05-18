package domain

import (
	"time"

	"github.com/google/uuid"
)

// HistoriqueAction represents the type of action logged.
type HistoriqueAction string

const (
	ActionCreation   HistoriqueAction = "CREATION"
	ActionUpdate     HistoriqueAction = "UPDATE"
	ActionStatus     HistoriqueAction = "STATUS_CHANGE"
	ActionAssignment HistoriqueAction = "ASSIGNMENT"
	ActionComment    HistoriqueAction = "COMMENT"
	ActionTimeLog    HistoriqueAction = "TIME_LOG"
)

// OrdreTravailHistorique logs every major change to an OrdreTravail.
type OrdreTravailHistorique struct {
	IDHistorique   uuid.UUID        `gorm:"column:id_historique;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_historique"`
	IDOrdreTravail uuid.UUID        `gorm:"column:id_ordre_travail;type:uuid;not null" json:"id_ordre_travail"`
	IDUtilisateur  uuid.UUID        `gorm:"column:id_utilisateur;type:uuid;not null" json:"id_utilisateur"`
	Action         HistoriqueAction `gorm:"column:action;type:varchar(30);not null" json:"action"`
	AncienStatut   *string          `gorm:"column:ancien_statut;type:varchar(50)" json:"ancien_statut"`
	NouveauStatut  *string          `gorm:"column:nouveau_statut;type:varchar(50)" json:"nouveau_statut"`
	Commentaire    string           `gorm:"column:commentaire" json:"commentaire"`
	CreatedAt      time.Time        `gorm:"column:created_at" json:"created_at"`
}

func (OrdreTravailHistorique) TableName() string {
	return "ordre_travail_historique"
}
