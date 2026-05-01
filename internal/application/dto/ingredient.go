package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateIngredientInput struct {
	Name         string
	Unit         string
	AverageCost  decimal.Decimal
	Stock        decimal.Decimal
	MinimumStock decimal.Decimal
}

type UpdateIngredientInput struct {
	ID           uuid.UUID
	Name         string
	Unit         string
	AverageCost  decimal.Decimal
	Stock        decimal.Decimal
	MinimumStock decimal.Decimal
}

type IngredientOutput struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Unit         string          `json:"unit"`
	AverageCost  decimal.Decimal `json:"average_cost"`
	Stock        decimal.Decimal `json:"stock"`
	MinimumStock decimal.Decimal `json:"minimum_stock"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	IsActive     bool            `json:"is_active"`
}
