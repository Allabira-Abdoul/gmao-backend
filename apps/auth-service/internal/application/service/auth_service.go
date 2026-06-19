package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend-gmao/apps/auth-service/internal/core/domain"
	"backend-gmao/apps/auth-service/internal/core/ports/secondary"
	"backend-gmao/pkg/audit"
	"backend-gmao/pkg/auth"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
)

type AuthService struct {
	sessionRepo secondary.SessionRepository
	// 🏛️ Atlas: Dependency Inversion Principle (DIP) applied.
	// AuthService now relies on the secondary port userClient instead of directly instantiating http.Client and using discovery.Registry
	// Decoupling from low-level communication logic reduces complexity and increases testability.
	userClient  secondary.UserClient
	jwtManager  *auth.JWTManager
	auditClient audit.Client
}

// NewAuthService creates a new authentication service.
func NewAuthService(
	sessionRepo secondary.SessionRepository,
	userClient secondary.UserClient,
	jwtManager *auth.JWTManager,
	auditClient audit.Client,
) *AuthService {
	return &AuthService{
		sessionRepo: sessionRepo,
		userClient:  userClient,
		jwtManager:  jwtManager,
		auditClient: auditClient,
	}
}

func (s *AuthService) CreateSession(ctx context.Context, req domain.CreateSessionRequest) (*domain.SessionResponse, error) {
	// 1. Fetch user by email via userClient
	user, err := s.userClient.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// 2. Verify user status
	if user.Status != domain.StatusActive {
		return nil, fmt.Errorf("account is %s", strings.ToLower(string(user.Status)))
	}

	// 4. Verify password
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// 5. Generate signed JWT access token
	jwtToken, expiredAt, err := s.jwtManager.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		user.RoleName,
		user.Privileges,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 6. Save session in DB
	session := &domain.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     jwtToken,
		ExpiredAt: expiredAt,
	}

	if _, err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	// Trigger audit event asynchronously
	userIDStr := user.ID.String()
	go func() {
		bgCtx := context.Background()
		_ = s.auditClient.LogEvent(bgCtx, audit.AuditEvent{
			ServiceName: "auth-service",
			Action:      "USER_LOGIN",
			Details:     fmt.Sprintf("User %s logged in successfully", user.Email),
			UserID:      &userIDStr,
		})
	}()

	sessionResp := session.ToResponse()
	return &sessionResp, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*domain.SessionResponse, error) {
	session, err := s.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiredAt) {
		s.sessionRepo.Logout(ctx, token)
		return nil, ErrSessionExpired
	}

	resp := session.ToResponse()
	return &resp, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, token string) error {
	return s.sessionRepo.Logout(ctx, token)
}
