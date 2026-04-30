package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend-gmao/apps/authentication-service/internal/core/domain"
	"backend-gmao/apps/authentication-service/internal/core/ports"
	"backend-gmao/pkg/auth"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountDisabled    = errors.New("account is disabled")
	ErrInvalidToken       = errors.New("invalid or expired refresh token")
)

type AuthService struct {
	tokenRepo  ports.TokenRepository
	userClient ports.UserClient
	jwtManager *auth.JWTManager
}

func NewAuthService(
	tokenRepo ports.TokenRepository,
	userClient ports.UserClient,
	jwtManager *auth.JWTManager,
) *AuthService {
	return &AuthService{
		tokenRepo:  tokenRepo,
		userClient: userClient,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.TokenPair, error) {
	// 1. Fetch user from user-service
	userInfo, err := s.userClient.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Verify password
	if !auth.CheckPasswordHash(req.MotDePasse, userInfo.MotDePasse) {
		return nil, ErrInvalidCredentials
	}

	// 3. Check account status
	if userInfo.StatutCompte != "ACTIVE" {
		return nil, ErrAccountDisabled
	}

	// 4. Generate token pair
	accessToken, _, err := s.jwtManager.GenerateAccessToken(
		userInfo.IDUtilisateur.String(),
		userInfo.Email,
		userInfo.RoleLibelle,
		userInfo.Privileges,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenStr, expiresAt, err := s.jwtManager.GenerateRefreshToken(userInfo.IDUtilisateur.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 5. Store refresh token
	tokenEntity := &domain.RefreshToken{
		UserID:    userInfo.IDUtilisateur,
		Token:     refreshTokenStr,
		ExpiresAt: expiresAt,
	}
	if err := s.tokenRepo.StoreRefreshToken(ctx, tokenEntity); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(s.jwtManager.GetAccessTokenDuration()),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req domain.RefreshRequest) (*domain.TokenPair, error) {
	// 1. Find token in DB
	storedToken, err := s.tokenRepo.FindRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 2. Validate token
	if storedToken.Revoked || time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// 3. Fetch latest user info
	userInfo, err := s.userClient.FindUserByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, ErrInvalidToken // Or specific error if user deleted
	}

	if userInfo.StatutCompte != "ACTIVE" {
		return nil, ErrAccountDisabled
	}

	// 4. Generate new pair (Rotate Refresh Token)
	accessToken, _, err := s.jwtManager.GenerateAccessToken(
		userInfo.IDUtilisateur.String(),
		userInfo.Email,
		userInfo.RoleLibelle,
		userInfo.Privileges,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshTokenStr, expiresAt, err := s.jwtManager.GenerateRefreshToken(userInfo.IDUtilisateur.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 5. Revoke old and store new (Rotation)
	if err := s.tokenRepo.RevokeRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	newEntity := &domain.RefreshToken{
		UserID:    userInfo.IDUtilisateur,
		Token:     newRefreshTokenStr,
		ExpiresAt: expiresAt,
	}
	if err := s.tokenRepo.StoreRefreshToken(ctx, newEntity); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenStr,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(s.jwtManager.GetAccessTokenDuration()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req domain.LogoutRequest) error {
	return s.tokenRepo.RevokeRefreshToken(ctx, req.RefreshToken)
}
