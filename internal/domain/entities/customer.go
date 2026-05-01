package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type Customer struct {
	ID           uuid.UUID
	FullName     string
	Phone        string
	Email        string
	CustomerType valueobjects.CustomerType
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsActive     bool
}
