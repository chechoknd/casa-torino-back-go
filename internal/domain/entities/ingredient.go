package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type Ingredient struct {
	ID           uuid.UUID
	Name         string
	Unit         valueobjects.Unit
	AverageCost  decimal.Decimal
	Stock        decimal.Decimal
	MinimumStock decimal.Decimal
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsActive     bool
}
