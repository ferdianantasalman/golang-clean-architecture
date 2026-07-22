package entity

import "github.com/google/uuid"

type RefreshToken struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey"`
	UserID    uuid.UUID `gorm:"column:user_id"`
	TokenHash string    `gorm:"column:token_hash"`
	ExpiresAt int64     `gorm:"column:expires_at"`
	CreatedAt int64     `gorm:"column:created_at;autoCreateTime:milli"`
}

func (r *RefreshToken) TableName() string {
	return "refresh_tokens"
}
