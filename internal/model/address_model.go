package model

import "github.com/google/uuid"

type AddressResponse struct {
	ID         uuid.UUID `json:"id"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	Province   string    `json:"province"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	CreatedAt  int64     `json:"created_at"`
	UpdatedAt  int64     `json:"updated_at"`
}

type ListAddressRequest struct {
	UserId    uuid.UUID `json:"-" validate:"required"`
	ContactId uuid.UUID `json:"-" validate:"required"`
}

type CreateAddressRequest struct {
	UserId     uuid.UUID `json:"-" validate:"required"`
	ContactId  uuid.UUID `json:"-" validate:"required"`
	Street     string    `json:"street" validate:"max=255"`
	City       string    `json:"city" validate:"max=255"`
	Province   string    `json:"province" validate:"max=255"`
	PostalCode string    `json:"postal_code" validate:"max=10"`
	Country    string    `json:"country" validate:"max=100"`
}

type UpdateAddressRequest struct {
	UserId     uuid.UUID `json:"-" validate:"required"`
	ContactId  uuid.UUID `json:"-" validate:"required"`
	ID         uuid.UUID `json:"-" validate:"required"`
	Street     string    `json:"street" validate:"max=255"`
	City       string    `json:"city" validate:"max=255"`
	Province   string    `json:"province" validate:"max=255"`
	PostalCode string    `json:"postal_code" validate:"max=10"`
	Country    string    `json:"country" validate:"max=100"`
}

type GetAddressRequest struct {
	UserId    uuid.UUID `json:"-" validate:"required"`
	ContactId uuid.UUID `json:"-" validate:"required"`
	ID        uuid.UUID `json:"-" validate:"required"`
}

type DeleteAddressRequest struct {
	UserId    uuid.UUID `json:"-" validate:"required"`
	ContactId uuid.UUID `json:"-" validate:"required"`
	ID        uuid.UUID `json:"-" validate:"required"`
}
