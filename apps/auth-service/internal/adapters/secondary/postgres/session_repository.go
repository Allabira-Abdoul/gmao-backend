package postgres

import (
	"context"
	"errors"
	"time"

	"backend-gmao/apps/auth-service/internal/core/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a GORM session repository.
func NewSessionRepository(db *gorm.DB) *sessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	err := r.db.WithContext(ctx).Create(session).Error
	return session, err
}

func (r *sessionRepository) Login(ctx context.Context, email string, password string) (*domain.Session, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("invalid password")
	}
	return r.Create(ctx, &domain.Session{
		UserID:    user.ID,
		Token:     uuid.New().String(),
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Hour * 24 * 7),
	})
}

func (r *sessionRepository) Logout(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&domain.Session{}, "token = ?", token).Error
}

func (r *sessionRepository) FindByToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).First(&session, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	var sessions []domain.Session
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}
