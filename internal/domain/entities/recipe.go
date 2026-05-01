package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type Recipe struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Name      string
	Portions  int
	Items     []RecipeItem
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}

func (r Recipe) CalculateCost(ingredients map[uuid.UUID]decimal.Decimal) (decimal.Decimal, error) {
	total := decimal.Zero

	for _, item := range r.Items {
		cost, ok := ingredients[item.IngredientID]
		if !ok {
			return decimal.Zero, fmt.Errorf("%w: ingredient cost %s", domainerrors.ErrNotFound, item.IngredientID)
		}

		total = total.Add(cost.Mul(item.Quantity))
	}

	return total, nil
}
