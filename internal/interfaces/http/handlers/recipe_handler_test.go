package handlers_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/application/dto"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/interfaces/http/handlers"
	"github.com/casatorino/backend/internal/interfaces/http/routes"
)

type fakeRecipeUseCase struct {
	createFn       func(context.Context, dto.CreateRecipeInput) (dto.RecipeOutput, error)
	addItemFn      func(context.Context, dto.AddRecipeItemInput) (dto.RecipeOutput, error)
	getByProductFn func(context.Context, uuid.UUID) (dto.RecipeOutput, error)
	listFn         func(context.Context) ([]dto.RecipeOutput, error)
	getCostFn      func(context.Context, uuid.UUID) (dto.RecipeCostOutput, error)
}

func (f fakeRecipeUseCase) CreateRecipe(ctx context.Context, input dto.CreateRecipeInput) (dto.RecipeOutput, error) {
	return f.createFn(ctx, input)
}

func (f fakeRecipeUseCase) AddRecipeItem(ctx context.Context, input dto.AddRecipeItemInput) (dto.RecipeOutput, error) {
	return f.addItemFn(ctx, input)
}

func (f fakeRecipeUseCase) GetRecipeByProduct(ctx context.Context, productID uuid.UUID) (dto.RecipeOutput, error) {
	return f.getByProductFn(ctx, productID)
}

func (f fakeRecipeUseCase) ListRecipes(ctx context.Context) ([]dto.RecipeOutput, error) {
	return f.listFn(ctx)
}

func (f fakeRecipeUseCase) GetRecipeCost(ctx context.Context, recipeID uuid.UUID) (dto.RecipeCostOutput, error) {
	return f.getCostFn(ctx, recipeID)
}

func TestRecipeEndpointsSuccess(t *testing.T) {
	recipeID := uuid.New()
	productID := uuid.New()
	ingredientID := uuid.New()
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	output := dto.RecipeOutput{
		ID:          recipeID,
		ProductID:   productID,
		ProductName: "Producto Demo",
		Name:        "Receta Demo",
		Portions:    1,
		Items: []dto.RecipeItemOutput{{
			ID:           uuid.New(),
			RecipeID:     recipeID,
			IngredientID: ingredientID,
			Quantity:     decimal.RequireFromString("0.25"),
			Unit:         "KG",
		}},
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}

	useCase := fakeRecipeUseCase{
		createFn: func(_ context.Context, input dto.CreateRecipeInput) (dto.RecipeOutput, error) {
			if input.ProductID != productID || input.Name != "Receta Demo" || input.Portions != 1 {
				t.Fatalf("unexpected create input: %+v", input)
			}
			created := output
			created.Items = nil
			return created, nil
		},
		addItemFn: func(_ context.Context, input dto.AddRecipeItemInput) (dto.RecipeOutput, error) {
			if input.RecipeID != recipeID || input.IngredientID != ingredientID || input.Unit != "KG" || !input.Quantity.Equal(decimal.RequireFromString("0.25")) {
				t.Fatalf("unexpected add item input: %+v", input)
			}
			return output, nil
		},
		getByProductFn: func(_ context.Context, id uuid.UUID) (dto.RecipeOutput, error) {
			if id != productID {
				t.Fatalf("unexpected product id: %s", id)
			}
			return output, nil
		},
		listFn: func(context.Context) ([]dto.RecipeOutput, error) {
			return []dto.RecipeOutput{output}, nil
		},
		getCostFn: func(_ context.Context, id uuid.UUID) (dto.RecipeCostOutput, error) {
			if id != recipeID {
				t.Fatalf("unexpected recipe id: %s", id)
			}
			return dto.RecipeCostOutput{RecipeID: recipeID, Cost: decimal.RequireFromString("1125")}, nil
		},
	}

	router := recipeTestRouter(useCase)
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"create", http.MethodPost, "/recipes/", `{"product_id":"` + productID.String() + `","name":"Receta Demo","portions":1}`, http.StatusCreated},
		{"list", http.MethodGet, "/recipes/", "", http.StatusOK},
		{"add item", http.MethodPost, "/recipes/" + recipeID.String() + "/items", `{"ingredient_id":"` + ingredientID.String() + `","quantity":"0.25","unit":"KG"}`, http.StatusOK},
		{"get by product", http.MethodGet, "/recipes/product/" + productID.String(), "", http.StatusOK},
		{"get cost", http.MethodGet, "/recipes/" + recipeID.String() + "/cost", "", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := performCustomerRequest(router, tc.method, tc.path, tc.body)
			if recorder.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body: %s", recorder.Code, tc.wantStatus, recorder.Body.String())
			}
			assertSuccessEnvelope(t, recorder.Body.Bytes())
		})
	}
}

func TestRecipeEndpointsErrors(t *testing.T) {
	recipeID := uuid.New()
	productID := uuid.New()
	useCase := fakeRecipeUseCase{
		createFn: func(context.Context, dto.CreateRecipeInput) (dto.RecipeOutput, error) {
			return dto.RecipeOutput{}, domainerrors.ErrInvalidInput
		},
		addItemFn: func(context.Context, dto.AddRecipeItemInput) (dto.RecipeOutput, error) {
			return dto.RecipeOutput{}, domainerrors.ErrInactive
		},
		getByProductFn: func(context.Context, uuid.UUID) (dto.RecipeOutput, error) {
			return dto.RecipeOutput{}, domainerrors.ErrNotFound
		},
		listFn: func(context.Context) ([]dto.RecipeOutput, error) {
			return nil, domainerrors.ErrInvalidInput
		},
		getCostFn: func(context.Context, uuid.UUID) (dto.RecipeCostOutput, error) {
			return dto.RecipeCostOutput{}, domainerrors.ErrNotFound
		},
	}

	router := recipeTestRouter(useCase)
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{"create invalid product uuid", http.MethodPost, "/recipes/", `{"product_id":"bad","name":"Receta Demo","portions":1}`, http.StatusBadRequest, "INVALID_INPUT"},
		{"list usecase error", http.MethodGet, "/recipes/", "", http.StatusBadRequest, "INVALID_INPUT"},
		{"add item invalid quantity", http.MethodPost, "/recipes/" + recipeID.String() + "/items", `{"ingredient_id":"` + productID.String() + `","quantity":"bad","unit":"KG"}`, http.StatusBadRequest, "INVALID_INPUT"},
		{"get by product not found", http.MethodGet, "/recipes/product/" + productID.String(), "", http.StatusNotFound, "NOT_FOUND"},
		{"get cost invalid uuid", http.MethodGet, "/recipes/not-a-uuid/cost", "", http.StatusBadRequest, "INVALID_INPUT"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := performCustomerRequest(router, tc.method, tc.path, tc.body)
			if recorder.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body: %s", recorder.Code, tc.wantStatus, recorder.Body.String())
			}
			assertErrorEnvelope(t, recorder.Body.Bytes(), tc.wantCode)
		})
	}
}

func recipeTestRouter(useCase fakeRecipeUseCase) http.Handler {
	return routes.NewRouter(routes.Dependencies{
		Customers:   handlers.NewCustomerHandler(nil),
		Products:    handlers.NewProductHandler(nil),
		Ingredients: handlers.NewIngredientHandler(nil),
		Recipes:     handlers.NewRecipeHandler(useCase),
		Orders:      handlers.NewOrderHandler(nil),
		Payments:    handlers.NewPaymentHandler(nil),
	})
}
