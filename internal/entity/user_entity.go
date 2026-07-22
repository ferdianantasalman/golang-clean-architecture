package entity

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey"`
	Password  string    `gorm:"column:password"`
	Name      string    `gorm:"column:name"`
	Token     string    `gorm:"column:token"`
	CreatedAt int64     `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64     `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Contacts  []Contact `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}
