package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/entities"
)

type IngredientRepository interface {
	Create(ctx context.Context, ingredient *entities.Ingredient) error
	Update(ctx context.Context, ingredient *entities.Ingredient) error
	Deactivate(ctx context.Context, id uuid.UUID, updatedAt time.Time) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Ingredient, error)
	ListActive(ctx context.Context) ([]entities.Ingredient, error)
}
