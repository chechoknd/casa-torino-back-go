package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateRecipeInput struct {
	ProductID uuid.UUID
	Name      string
	Portions  int
}

type AddRecipeItemInput struct {
	RecipeID     uuid.UUID
	IngredientID uuid.UUID
	Quantity     decimal.Decimal
	Unit         string
}

type RecipeItemOutput struct {
	ID           uuid.UUID       `json:"id"`
	RecipeID     uuid.UUID       `json:"recipe_id"`
	IngredientID uuid.UUID       `json:"ingredient_id"`
	Quantity     decimal.Decimal `json:"quantity"`
	Unit         string          `json:"unit"`
}

type RecipeOutput struct {
	ID          uuid.UUID          `json:"id"`
	ProductID   uuid.UUID          `json:"product_id"`
	ProductName string             `json:"product_name"`
	Name        string             `json:"name"`
	Portions    int                `json:"portions"`
	Items       []RecipeItemOutput `json:"items"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	IsActive    bool               `json:"is_active"`
}

type RecipeCostOutput struct {
	RecipeID uuid.UUID       `json:"recipe_id"`
	Cost     decimal.Decimal `json:"cost"`
}
