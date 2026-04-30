package postgres

import (
	"context"
	"time"

	"backend-gmao/apps/authentication-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *tokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) StoreRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) FindRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *tokenRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

func (r *tokenRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

func (r *tokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ? OR revoked = ?", time.Now(), true).
		Delete(&domain.RefreshToken{}).Error
}
