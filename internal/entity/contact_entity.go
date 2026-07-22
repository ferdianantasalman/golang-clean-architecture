package entity

import "github.com/google/uuid"

type Contact struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Email     string    `gorm:"column:email"`
	Phone     string    `gorm:"column:phone"`
	UserId    uuid.UUID `gorm:"column:user_id"`
	CreatedAt int64     `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64     `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	User      User      `gorm:"foreignKey:user_id;references:id"`
	Addresses []Address `gorm:"foreignKey:contact_id;references:id"`
}

func (c *Contact) TableName() string {
	return "contacts"
}
