package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCustomerInput struct {
	FullName     string
	Phone        string
	Email        string
	CustomerType string
}

type UpdateCustomerInput struct {
	ID           uuid.UUID
	FullName     string
	Phone        string
	Email        string
	CustomerType string
}

type CustomerOutput struct {
	ID           uuid.UUID `json:"id"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	CustomerType string    `json:"customer_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsActive     bool      `json:"is_active"`
}
