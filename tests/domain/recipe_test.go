package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/domain/entities"
)

func TestRecipeCalculateCost(t *testing.T) {
	t.Parallel()

	ingredientA := uuid.New()
	ingredientB := uuid.New()

	recipe := entities.Recipe{
		Items: []entities.RecipeItem{
			{IngredientID: ingredientA, Quantity: decimal.RequireFromString("0.5")},
			{IngredientID: ingredientB, Quantity: decimal.RequireFromString("2")},
		},
	}

	cost, err := recipe.CalculateCost(map[uuid.UUID]decimal.Decimal{
		ingredientA: decimal.RequireFromString("10000"),
		ingredientB: decimal.RequireFromString("1500"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cost.Equal(decimal.RequireFromString("8000")) {
		t.Fatalf("unexpected recipe cost: %s", cost)
	}
}
