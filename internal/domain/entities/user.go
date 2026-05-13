package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	FullName     string
	Role         valueobjects.UserRole
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsActive     bool
}
