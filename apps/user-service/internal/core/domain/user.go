package domain

import (
	"time"

	"github.com/google/uuid"
)

// AccountStatus represents the status of a user account.
type AccountStatus string

const (
	StatusActive   AccountStatus = "ACTIVE"
	StatusInactive AccountStatus = "INACTIVE"
	StatusLocked   AccountStatus = "LOCKED"
)

// User represents the Utilisateur entity in the GMAO system.
type User struct {
	IDUtilisateur uuid.UUID     `gorm:"column:id_utilisateur;type:uuid;primaryKey;default:gen_random_uuid()" json:"id_utilisateur"`
	NomComplet    string        `gorm:"column:nom_complet;not null" json:"nom_complet"`
	Email         string        `gorm:"column:email;uniqueIndex;not null" json:"email"`
	MotDePasse    string        `gorm:"column:mot_de_passe;not null" json:"-"`
	StatutCompte  AccountStatus `gorm:"column:statut_compte;type:varchar(20);default:'ACTIVE'" json:"statut_compte"`
	IDRole        uuid.UUID     `gorm:"column:id_role;type:uuid;not null" json:"id_role"`
	Role          Role          `gorm:"foreignKey:IDRole;references:IDRole" json:"role,omitempty"`
	CreatedAt     time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "utilisateurs"
}

// UserResponse is the DTO returned by API endpoints (excludes password).
type UserResponse struct {
	IDUtilisateur uuid.UUID     `json:"id_utilisateur"`
	NomComplet    string        `json:"nom_complet"`
	Email         string        `json:"email"`
	StatutCompte  AccountStatus `json:"statut_compte"`
	IDRole        uuid.UUID     `json:"id_role"`
	Role          *RoleResponse `json:"role,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// ToResponse converts a User to a UserResponse (safe for API output).
func (u *User) ToResponse() UserResponse {
	resp := UserResponse{
		IDUtilisateur: u.IDUtilisateur,
		NomComplet:    u.NomComplet,
		Email:         u.Email,
		StatutCompte:  u.StatutCompte,
		IDRole:        u.IDRole,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}

	if u.Role.IDRole != uuid.Nil {
		roleResp := u.Role.ToResponse()
		resp.Role = &roleResp
	}

	return resp
}

// InternalUserResponse is the DTO used for inter-service communication.
// It includes the hashed password for authentication verification.
type InternalUserResponse struct {
	IDUtilisateur uuid.UUID     `json:"id_utilisateur"`
	NomComplet    string        `json:"nom_complet"`
	Email         string        `json:"email"`
	MotDePasse    string        `json:"mot_de_passe"`
	StatutCompte  AccountStatus `json:"statut_compte"`
	RoleLibelle   string        `json:"role_libelle"`
	Privileges    []string      `json:"privileges"`
}

// ToInternalResponse converts a User to an InternalUserResponse.
func (u *User) ToInternalResponse() InternalUserResponse {
	privileges := make([]string, 0)
	for _, rp := range u.Role.Privileges {
		privileges = append(privileges, rp.Privilege)
	}

	return InternalUserResponse{
		IDUtilisateur: u.IDUtilisateur,
		NomComplet:    u.NomComplet,
		Email:         u.Email,
		MotDePasse:    u.MotDePasse,
		StatutCompte:  u.StatutCompte,
		RoleLibelle:   u.Role.Libelle,
		Privileges:    privileges,
	}
}

// CreateUserRequest is the DTO for creating a new user.
type CreateUserRequest struct {
	NomComplet string `json:"nom_complet" binding:"required,min=2,max=255"`
	Email      string `json:"email" binding:"required,email"`
	MotDePasse string `json:"mot_de_passe" binding:"required,min=8"`
	IDRole     string `json:"id_role" binding:"required,uuid"`
}

// UpdateUserRequest is the DTO for updating an existing user.
type UpdateUserRequest struct {
	NomComplet   *string `json:"nom_complet,omitempty" binding:"omitempty,min=2,max=255"`
	Email        *string `json:"email,omitempty" binding:"omitempty,email"`
	StatutCompte *string `json:"statut_compte,omitempty" binding:"omitempty,oneof=ACTIVE INACTIVE LOCKED"`
	IDRole       *string `json:"id_role,omitempty" binding:"omitempty,uuid"`
}
