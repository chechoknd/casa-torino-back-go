package entities

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type RecipeItem struct {
	ID           uuid.UUID
	RecipeID     uuid.UUID
	IngredientID uuid.UUID
	Quantity     decimal.Decimal
	Unit         valueobjects.Unit
}
