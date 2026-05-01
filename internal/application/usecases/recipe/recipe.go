package recipe

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/repositories"
	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type UseCase struct {
	recipes     repositories.RecipeRepository
	products    repositories.ProductRepository
	ingredients repositories.IngredientRepository
}

func NewUseCase(recipes repositories.RecipeRepository, products repositories.ProductRepository, ingredients repositories.IngredientRepository) *UseCase {
	return &UseCase{
		recipes:     recipes,
		products:    products,
		ingredients: ingredients,
	}
}

func (uc *UseCase) CreateRecipe(ctx context.Context, input dto.CreateRecipeInput) (dto.RecipeOutput, error) {
	if strings.TrimSpace(input.Name) == "" || input.Portions <= 0 {
		return dto.RecipeOutput{}, domainerrors.ErrInvalidInput
	}

	product, err := uc.products.FindByID(ctx, input.ProductID)
	if err != nil {
		return dto.RecipeOutput{}, err
	}
	if !product.IsActive {
		return dto.RecipeOutput{}, domainerrors.ErrInactive
	}

	now := time.Now().UTC()
	recipe := &entities.Recipe{
		ID:        uuid.New(),
		ProductID: input.ProductID,
		Name:      strings.TrimSpace(input.Name),
		Portions:  input.Portions,
		Items:     []entities.RecipeItem{},
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}

	if err := uc.recipes.Create(ctx, recipe); err != nil {
		return dto.RecipeOutput{}, err
	}

	return uc.toRecipeOutput(ctx, *recipe)
}

func (uc *UseCase) AddRecipeItem(ctx context.Context, input dto.AddRecipeItemInput) (dto.RecipeOutput, error) {
	if input.Quantity.LessThanOrEqual(decimal.Zero) {
		return dto.RecipeOutput{}, domainerrors.ErrInvalidInput
	}

	recipe, err := uc.recipes.FindByID(ctx, input.RecipeID)
	if err != nil {
		return dto.RecipeOutput{}, err
	}
	if !recipe.IsActive {
		return dto.RecipeOutput{}, domainerrors.ErrInactive
	}

	ingredient, err := uc.ingredients.FindByID(ctx, input.IngredientID)
	if err != nil {
		return dto.RecipeOutput{}, err
	}
	if !ingredient.IsActive {
		return dto.RecipeOutput{}, domainerrors.ErrInactive
	}

	unit, err := valueobjects.NewUnit(input.Unit)
	if err != nil {
		return dto.RecipeOutput{}, err
	}

	item := &entities.RecipeItem{
		ID:           uuid.New(),
		RecipeID:     recipe.ID,
		IngredientID: input.IngredientID,
		Quantity:     input.Quantity,
		Unit:         unit,
	}

	if err := uc.recipes.AddItem(ctx, recipe.ID, item); err != nil {
		return dto.RecipeOutput{}, err
	}

	updatedRecipe, err := uc.recipes.FindByID(ctx, recipe.ID)
	if err != nil {
		return dto.RecipeOutput{}, err
	}

	return uc.toRecipeOutput(ctx, *updatedRecipe)
}

func (uc *UseCase) GetRecipeByProduct(ctx context.Context, productID uuid.UUID) (dto.RecipeOutput, error) {
	recipe, err := uc.recipes.FindByProductID(ctx, productID)
	if err != nil {
		return dto.RecipeOutput{}, err
	}
	if !recipe.IsActive {
		return dto.RecipeOutput{}, domainerrors.ErrInactive
	}

	return uc.toRecipeOutput(ctx, *recipe)
}

func (uc *UseCase) ListRecipes(ctx context.Context) ([]dto.RecipeOutput, error) {
	recipes, err := uc.recipes.List(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]dto.RecipeOutput, 0, len(recipes))
	for _, recipe := range recipes {
		mapped, err := uc.toRecipeOutput(ctx, recipe)
		if err != nil {
			return nil, err
		}
		output = append(output, mapped)
	}

	return output, nil
}

func (uc *UseCase) CalculateRecipeCost(ctx context.Context, productID uuid.UUID) (dto.RecipeCostOutput, error) {
	recipe, err := uc.recipes.FindByProductID(ctx, productID)
	if err != nil {
		return dto.RecipeCostOutput{}, err
	}

	return uc.calculateRecipeCost(ctx, recipe)
}

func (uc *UseCase) GetRecipeCost(ctx context.Context, recipeID uuid.UUID) (dto.RecipeCostOutput, error) {
	recipe, err := uc.recipes.FindByID(ctx, recipeID)
	if err != nil {
		return dto.RecipeCostOutput{}, err
	}

	return uc.calculateRecipeCost(ctx, recipe)
}

func (uc *UseCase) calculateRecipeCost(ctx context.Context, recipe *entities.Recipe) (dto.RecipeCostOutput, error) {
	costs := make(map[uuid.UUID]decimal.Decimal, len(recipe.Items))
	for _, item := range recipe.Items {
		ingredient, err := uc.ingredients.FindByID(ctx, item.IngredientID)
		if err != nil {
			return dto.RecipeCostOutput{}, err
		}
		costs[item.IngredientID] = ingredient.AverageCost
	}

	total, err := recipe.CalculateCost(costs)
	if err != nil {
		return dto.RecipeCostOutput{}, err
	}

	return dto.RecipeCostOutput{
		RecipeID: recipe.ID,
		Cost:     total,
	}, nil
}

func (uc *UseCase) toRecipeOutput(ctx context.Context, recipe entities.Recipe) (dto.RecipeOutput, error) {
	productName := ""
	product, err := uc.products.FindByID(ctx, recipe.ProductID)
	if err == nil {
		productName = product.Name
	}

	items := make([]dto.RecipeItemOutput, 0, len(recipe.Items))
	for _, item := range recipe.Items {
		items = append(items, dto.RecipeItemOutput{
			ID:           item.ID,
			RecipeID:     item.RecipeID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			Unit:         string(item.Unit),
		})
	}

	return dto.RecipeOutput{
		ID:          recipe.ID,
		ProductID:   recipe.ProductID,
		ProductName: productName,
		Name:        recipe.Name,
		Portions:    recipe.Portions,
		Items:       items,
		CreatedAt:   recipe.CreatedAt,
		UpdatedAt:   recipe.UpdatedAt,
		IsActive:    recipe.IsActive,
	}, nil
}
