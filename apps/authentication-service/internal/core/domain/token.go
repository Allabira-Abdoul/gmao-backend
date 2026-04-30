package domain

import (
	"time"

	"github.com/google/uuid"
)

// TokenPair represents an access + refresh token pair returned after authentication.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshToken represents a stored refresh token in the database.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"column:token;uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	Revoked   bool      `gorm:"column:revoked;default:false" json:"revoked"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName overrides the default table name.
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// LoginRequest is the DTO for login authentication.
type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	MotDePasse string `json:"mot_de_passe" binding:"required"`
}

// RefreshRequest is the DTO for token refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest is the DTO for logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UserInfo represents the user data received from user-service for authentication.
type UserInfo struct {
	IDUtilisateur uuid.UUID `json:"id_utilisateur"`
	NomComplet    string    `json:"nom_complet"`
	Email         string    `json:"email"`
	MotDePasse    string    `json:"mot_de_passe"`
	StatutCompte  string    `json:"statut_compte"`
	RoleLibelle   string    `json:"role_libelle"`
	Privileges    []string  `json:"privileges"`
}
