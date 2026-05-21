package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"backend-gmao/apps/auth-service/internal/core/domain"
	"backend-gmao/apps/auth-service/internal/core/ports/secondary"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
)

// AuthService implements primary.AuthService.
type AuthService struct {
	sessionRepo secondary.SessionRepository
}

// NewAuthService initializes a new AuthService instance.
func NewAuthService(sessionRepo secondary.SessionRepository) *AuthService {
	return &AuthService{sessionRepo: sessionRepo}
}

func (s *AuthService) CreateSession(ctx context.Context, req domain.CreateSessionRequest) (*domain.SessionResponse, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	session := &domain.Session{
		ID:        uuid.New(),
		UserID:    req.UserID,
		Token:     token,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	resp := session.ToResponse()
	return &resp, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*domain.SessionResponse, error) {
	session, err := s.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiredAt) {
		s.sessionRepo.Delete(ctx, token)
		return nil, ErrSessionExpired
	}

	resp := session.ToResponse()
	return &resp, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, token string) error {
	return s.sessionRepo.Delete(ctx, token)
}
