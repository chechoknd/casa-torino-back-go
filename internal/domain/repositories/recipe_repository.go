package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/entities"
)

type RecipeRepository interface {
	Create(ctx context.Context, recipe *entities.Recipe) error
	Update(ctx context.Context, recipe *entities.Recipe) error
	AddItem(ctx context.Context, recipeID uuid.UUID, item *entities.RecipeItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Recipe, error)
	FindByProductID(ctx context.Context, productID uuid.UUID) (*entities.Recipe, error)
}
