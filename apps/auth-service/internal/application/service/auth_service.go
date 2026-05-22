package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"backend-gmao/apps/auth-service/internal/core/domain"
	"backend-gmao/apps/auth-service/internal/core/ports/secondary"
	"backend-gmao/pkg/audit"
	"backend-gmao/pkg/auth"
	"backend-gmao/pkg/discovery"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
)

type AuthService struct {
	sessionRepo secondary.SessionRepository
	registry    discovery.Registry
	jwtManager  *auth.JWTManager
	auditClient audit.Client
}

// NewAuthService creates a new authentication service.
func NewAuthService(
	sessionRepo secondary.SessionRepository,
	registry discovery.Registry,
	jwtManager *auth.JWTManager,
	auditClient audit.Client,
) *AuthService {
	return &AuthService{
		sessionRepo: sessionRepo,
		registry:    registry,
		jwtManager:  jwtManager,
		auditClient: auditClient,
	}
}

func (s *AuthService) CreateSession(ctx context.Context, req domain.CreateSessionRequest) (*domain.SessionResponse, error) {
	// 1. Discover user-service via Consul
	addr, err := s.registry.Discover("user-service")
	if err != nil {
		return nil, fmt.Errorf("failed to discover user-service: %w", err)
	}

	// 2. Query user-service/internal/by-email
	client := &http.Client{Timeout: 5 * time.Second}
	targetURL := fmt.Sprintf("http://%s/internal/by-email?email=%s", addr, url.QueryEscape(req.Email))

	httpReq, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("X-Internal-Service", "auth-service")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("invalid email or password")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status: %d", resp.StatusCode)
	}

	var envelope struct {
		Status string      `json:"status"`
		Data   domain.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	user := envelope.Data

	// 3. Verify user status
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

