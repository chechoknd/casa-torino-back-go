package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casatorino/backend/internal/domain/entities"
	domainrepositories "github.com/casatorino/backend/internal/domain/repositories"
	sqlcdb "github.com/casatorino/backend/internal/infrastructure/database/sqlc"
)

var _ domainrepositories.IngredientRepository = (*IngredientRepository)(nil)

type IngredientRepository struct {
	conn    *pgxpool.Pool
	queries *sqlcdb.Queries
}

func NewIngredientRepository(conn *pgxpool.Pool) *IngredientRepository {
	return &IngredientRepository{
		conn:    conn,
		queries: sqlcdb.New(conn),
	}
}

func (r *IngredientRepository) Create(ctx context.Context, ingredient *entities.Ingredient) error {
	row, err := r.queries.CreateIngredient(ctx, sqlcdb.CreateIngredientParams{
		ID:           ingredient.ID,
		Name:         ingredient.Name,
		Unit:         string(ingredient.Unit),
		AverageCost:  ingredient.AverageCost,
		Stock:        ingredient.Stock,
		MinimumStock: ingredient.MinimumStock,
		CreatedAt:    ingredient.CreatedAt,
		UpdatedAt:    ingredient.UpdatedAt,
		IsActive:     ingredient.IsActive,
	})
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapIngredient(row)
	if err != nil {
		return err
	}

	*ingredient = mapped
	return nil
}

func (r *IngredientRepository) Update(ctx context.Context, ingredient *entities.Ingredient) error {
	row, err := r.queries.UpdateIngredient(ctx, sqlcdb.UpdateIngredientParams{
		ID:           ingredient.ID,
		Name:         ingredient.Name,
		Unit:         string(ingredient.Unit),
		AverageCost:  ingredient.AverageCost,
		Stock:        ingredient.Stock,
		MinimumStock: ingredient.MinimumStock,
		UpdatedAt:    ingredient.UpdatedAt,
		IsActive:     ingredient.IsActive,
	})
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapIngredient(row)
	if err != nil {
		return err
	}

	*ingredient = mapped
	return nil
}

func (r *IngredientRepository) Deactivate(ctx context.Context, id uuid.UUID, updatedAt time.Time) error {
	rows, err := r.queries.DeactivateIngredient(ctx, sqlcdb.DeactivateIngredientParams{
		ID:        id,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return mapError(err)
	}
	if rows == 0 {
		return domainrepositoriesErrNotFound()
	}
	return nil
}

func (r *IngredientRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Ingredient, error) {
	row, err := r.queries.GetIngredientByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	mapped, err := mapIngredient(row)
	if err != nil {
		return nil, err
	}

	return &mapped, nil
}

func (r *IngredientRepository) ListActive(ctx context.Context) ([]entities.Ingredient, error) {
	rows, err := r.queries.ListIngredients(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	ingredients := make([]entities.Ingredient, 0, len(rows))
	for _, row := range rows {
		ingredient, err := mapIngredient(row)
		if err != nil {
			return nil, err
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}
