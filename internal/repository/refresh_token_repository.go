package repository

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	Repository[entity.RefreshToken]
	Log *logrus.Logger
}

func NewRefreshTokenRepository(log *logrus.Logger) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		Log: log,
	}
}

func (r *RefreshTokenRepository) FindByHash(db *gorm.DB, refreshToken *entity.RefreshToken, hash string) error {
	return db.Where("token_hash = ?", hash).First(refreshToken).Error
}

func (r *RefreshTokenRepository) DeleteByUserId(db *gorm.DB, userId uuid.UUID) error {
	return db.Where("user_id = ?", userId).Delete(&entity.RefreshToken{}).Error
}
