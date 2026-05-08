package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	FullName     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsActive     bool
}
