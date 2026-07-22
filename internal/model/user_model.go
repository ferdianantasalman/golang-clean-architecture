package model

import "github.com/google/uuid"

type UserResponse struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	CreatedAt    int64     `json:"created_at,omitempty"`
	UpdatedAt    int64     `json:"updated_at,omitempty"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=1024"`
}

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID `json:"-" validate:"required"`
	Name     string    `json:"name,omitempty" validate:"max=100"`
	Password string    `json:"password,omitempty" validate:"max=100"`
}

type LoginUserRequest struct {
	Name     string `json:"name,omitempty" validate:"max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LogoutUserRequest struct {
	ID    uuid.UUID `json:"id" validate:"required"`
	Token string    `json:"-"`
}

type GetUserRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
