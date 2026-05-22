package domain

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user session in the auth system.
type Session struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	Token     string    `gorm:"column:token;uniqueIndex;not null" json:"token"`
	ExpiredAt time.Time `gorm:"column:expired_at;not null" json:"expired_at"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName overrides GORM's default table name.
func (Session) TableName() string {
	return "sessions"
}

// SessionResponse is the DTO returned by API endpoints.
type SessionResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a Session to SessionResponse DTO.
func (s *Session) ToResponse() SessionResponse {
	return SessionResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		Token:     s.Token,
		ExpiredAt: s.ExpiredAt,
		CreatedAt: s.CreatedAt,
	}
}

// CreateSessionRequest is the DTO used to submit a new session.
type CreateSessionRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountStatus string

const (
	StatusActive   AccountStatus = "ACTIVE"
	StatusInactive AccountStatus = "INACTIVE"
	StatusLocked   AccountStatus = "LOCKED"
)

type User struct {
	ID            uuid.UUID     `json:"id"`
	FullName      string        `json:"full_name"`
	Email         string        `json:"email"`
	Password      string        `json:"password"`
	Status        AccountStatus `json:"status"`
	RoleName      string        `json:"role_name"`
	Privileges    []string      `json:"privileges"`
}
