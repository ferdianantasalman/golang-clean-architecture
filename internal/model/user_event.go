package model

import "github.com/google/uuid"

type UserEvent struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt int64     `json:"created_at,omitempty"`
	UpdatedAt int64     `json:"updated_at,omitempty"`
}

func (u *UserEvent) GetId() uuid.UUID {
	return u.ID
}
