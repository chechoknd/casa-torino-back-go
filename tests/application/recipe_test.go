package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/application/dto"
	recipeuc "github.com/casatorino/backend/internal/application/usecases/recipe"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	"github.com/casatorino/backend/tests/mocks"
)

func TestCreateRecipeSuccess(t *testing.T) {
	pt, _ := valueobjects.NewProductType("LUNCH")
	productRepo := &mocks.ProductRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Product, error) {
			return &entities.Product{ID: uuid.New(), ProductType: pt, IsActive: true}, nil
		},
	}
	recipeRepo := &mocks.RecipeRepository{CreateFn: func(context.Context, *entities.Recipe) error { return nil }}
	ingredientRepo := &mocks.IngredientRepository{}
	uc := recipeuc.NewUseCase(recipeRepo, productRepo, ingredientRepo)

	_, err := uc.CreateRecipe(context.Background(), dto.CreateRecipeInput{ProductID: uuid.New(), Name: "Receta", Portions: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddRecipeItemInvalidQuantity(t *testing.T) {
	uc := recipeuc.NewUseCase(&mocks.RecipeRepository{}, &mocks.ProductRepository{}, &mocks.IngredientRepository{})
	_, err := uc.AddRecipeItem(context.Background(), dto.AddRecipeItemInput{RecipeID: uuid.New(), IngredientID: uuid.New(), Quantity: decimal.Zero, Unit: "G"})
	if !errors.Is(err, domainerrors.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestCalculateRecipeCostSuccess(t *testing.T) {
	recipeID := uuid.New()
	productID := uuid.New()
	ingredientID := uuid.New()
	unit, _ := valueobjects.NewUnit("G")
	recipeRepo := &mocks.RecipeRepository{
		FindByProductIDFn: func(context.Context, uuid.UUID) (*entities.Recipe, error) {
			return &entities.Recipe{
				ID: recipeID, ProductID: productID, IsActive: true,
				Items: []entities.RecipeItem{{IngredientID: ingredientID, Quantity: decimal.RequireFromString("2"), Unit: unit}},
			}, nil
		},
	}
	ingredientRepo := &mocks.IngredientRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Ingredient, error) {
			return &entities.Ingredient{ID: ingredientID, AverageCost: decimal.RequireFromString("3500"), IsActive: true}, nil
		},
	}
	uc := recipeuc.NewUseCase(recipeRepo, &mocks.ProductRepository{}, ingredientRepo)

	out, err := uc.CalculateRecipeCost(context.Background(), productID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Cost.Equal(decimal.RequireFromString("7000")) {
		t.Fatalf("unexpected cost: %s", out.Cost)
	}
}
