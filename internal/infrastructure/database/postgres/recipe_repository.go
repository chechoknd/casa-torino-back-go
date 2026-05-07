package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casatorino/backend/internal/domain/entities"
	domainrepositories "github.com/casatorino/backend/internal/domain/repositories"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	sqlcdb "github.com/casatorino/backend/internal/infrastructure/database/sqlc"
)

var _ domainrepositories.RecipeRepository = (*RecipeRepository)(nil)

type RecipeRepository struct {
	conn    *pgxpool.Pool
	queries *sqlcdb.Queries
}

func NewRecipeRepository(conn *pgxpool.Pool) *RecipeRepository {
	return &RecipeRepository{
		conn:    conn,
		queries: sqlcdb.New(conn),
	}
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *entities.Recipe) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)
	row, err := qtx.CreateRecipe(ctx, sqlcdb.CreateRecipeParams{
		ID:        recipe.ID,
		ProductID: recipe.ProductID,
		Name:      recipe.Name,
		Portions:  int32(recipe.Portions),
		CreatedAt: recipe.CreatedAt,
		UpdatedAt: recipe.UpdatedAt,
		IsActive:  recipe.IsActive,
	})
	if err != nil {
		return mapError(err)
	}

	for _, item := range recipe.Items {
		if _, err := qtx.CreateRecipeItem(ctx, sqlcdb.CreateRecipeItemParams{
			ID:           item.ID,
			RecipeID:     recipe.ID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			Unit:         string(item.Unit),
		}); err != nil {
			return mapError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	items, err := r.queries.GetRecipeItemsByRecipeID(ctx, row.ID)
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapRecipe(row, items)
	if err != nil {
		return err
	}

	*recipe = mapped
	return nil
}

func (r *RecipeRepository) Update(ctx context.Context, recipe *entities.Recipe) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)
	row, err := qtx.UpdateRecipe(ctx, sqlcdb.UpdateRecipeParams{
		ID:        recipe.ID,
		ProductID: recipe.ProductID,
		Name:      recipe.Name,
		Portions:  int32(recipe.Portions),
		UpdatedAt: recipe.UpdatedAt,
		IsActive:  recipe.IsActive,
	})
	if err != nil {
		return mapError(err)
	}

	for _, item := range recipe.Items {
		if item.ID == uuid.Nil {
			continue
		}

		if _, err := qtx.UpdateRecipeItem(ctx, sqlcdb.UpdateRecipeItemParams{
			ID:           item.ID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			Unit:         string(item.Unit),
		}); err != nil {
			return mapError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	items, err := r.queries.GetRecipeItemsByRecipeID(ctx, row.ID)
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapRecipe(row, items)
	if err != nil {
		return err
	}

	*recipe = mapped
	return nil
}

func (r *RecipeRepository) AddItem(ctx context.Context, recipeID uuid.UUID, item *entities.RecipeItem) error {
	row, err := r.queries.CreateRecipeItem(ctx, sqlcdb.CreateRecipeItemParams{
		ID:           item.ID,
		RecipeID:     recipeID,
		IngredientID: item.IngredientID,
		Quantity:     item.Quantity,
		Unit:         string(item.Unit),
	})
	if err != nil {
		return mapError(err)
	}

	unit, err := valueobjects.NewUnit(row.Unit)
	if err != nil {
		return err
	}

	item.ID = row.ID
	item.RecipeID = row.RecipeID
	item.IngredientID = row.IngredientID
	item.Quantity = row.Quantity
	item.Unit = unit

	return nil
}

func (r *RecipeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Recipe, error) {
	row, err := r.queries.GetRecipeByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	items, err := r.queries.GetRecipeItemsByRecipeID(ctx, row.ID)
	if err != nil {
		return nil, mapError(err)
	}

	mapped, err := mapRecipe(row, items)
	if err != nil {
		return nil, err
	}

	return &mapped, nil
}

func (r *RecipeRepository) FindByProductID(ctx context.Context, productID uuid.UUID) (*entities.Recipe, error) {
	row, err := r.queries.GetRecipeByProductID(ctx, productID)
	if err != nil {
		return nil, mapError(err)
	}

	items, err := r.queries.GetRecipeItemsByRecipeID(ctx, row.ID)
	if err != nil {
		return nil, mapError(err)
	}

	mapped, err := mapRecipe(row, items)
	if err != nil {
		return nil, err
	}

	return &mapped, nil
}

func (r *RecipeRepository) List(ctx context.Context) ([]entities.Recipe, error) {
	rows, err := r.queries.ListRecipes(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	recipes := make([]entities.Recipe, 0, len(rows))
	for _, row := range rows {
		items, err := r.queries.GetRecipeItemsByRecipeID(ctx, row.ID)
		if err != nil {
			return nil, mapError(err)
		}

		recipe, err := mapRecipe(row, items)
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
