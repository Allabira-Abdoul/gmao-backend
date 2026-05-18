package domain

import (
	"time"

	"github.com/google/uuid"
)

// WorkflowState represents a possible status in a workflow (e.g., PLANIFIE, EN_COURS).
type WorkflowState struct {
	IDState     uuid.UUID `gorm:"column:id_state;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_state"`
	Nom         string    `gorm:"column:nom;uniqueIndex;not null" json:"nom"` // e.g., "PLANIFIE"
	Description string    `gorm:"column:description" json:"description"`
	IsInitial   bool      `gorm:"column:is_initial;default:false" json:"is_initial"`
	IsFinal     bool      `gorm:"column:is_final;default:false" json:"is_final"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (WorkflowState) TableName() string {
	return "workflow_states"
}

// WorkflowTransition represents an allowed transition between two WorkflowStates.
type WorkflowTransition struct {
	IDTransition uuid.UUID `gorm:"column:id_transition;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_transition"`
	FromStateID  uuid.UUID `gorm:"column:from_state_id;type:uuid;not null" json:"from_state_id"`
	ToStateID    uuid.UUID `gorm:"column:to_state_id;type:uuid;not null" json:"to_state_id"`
	RequiredRole string    `gorm:"column:required_role" json:"required_role"` // e.g., "ADMIN", "TECHNICIEN"
	ActionName   string    `gorm:"column:action_name;not null" json:"action_name"` // e.g., "Commencer", "Terminer"
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (WorkflowTransition) TableName() string {
	return "workflow_transitions"
}
