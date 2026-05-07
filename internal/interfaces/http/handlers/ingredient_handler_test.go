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

type fakeIngredientUseCase struct {
	createFn     func(context.Context, dto.CreateIngredientInput) (dto.IngredientOutput, error)
	getFn        func(context.Context, uuid.UUID) (dto.IngredientOutput, error)
	listFn       func(context.Context) ([]dto.IngredientOutput, error)
	updateFn     func(context.Context, dto.UpdateIngredientInput) (dto.IngredientOutput, error)
	deactivateFn func(context.Context, uuid.UUID) error
}

func (f fakeIngredientUseCase) CreateIngredient(ctx context.Context, input dto.CreateIngredientInput) (dto.IngredientOutput, error) {
	return f.createFn(ctx, input)
}

func (f fakeIngredientUseCase) GetIngredient(ctx context.Context, id uuid.UUID) (dto.IngredientOutput, error) {
	return f.getFn(ctx, id)
}

func (f fakeIngredientUseCase) ListIngredients(ctx context.Context) ([]dto.IngredientOutput, error) {
	return f.listFn(ctx)
}

func (f fakeIngredientUseCase) UpdateIngredient(ctx context.Context, input dto.UpdateIngredientInput) (dto.IngredientOutput, error) {
	return f.updateFn(ctx, input)
}

func (f fakeIngredientUseCase) DeactivateIngredient(ctx context.Context, id uuid.UUID) error {
	return f.deactivateFn(ctx, id)
}

func TestIngredientEndpointsSuccess(t *testing.T) {
	ingredientID := uuid.New()
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	output := dto.IngredientOutput{
		ID:           ingredientID,
		Name:         "Ingrediente Demo",
		Unit:         "KG",
		AverageCost:  decimal.RequireFromString("4500"),
		Stock:        decimal.RequireFromString("20"),
		MinimumStock: decimal.RequireFromString("5"),
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}

	useCase := fakeIngredientUseCase{
		createFn: func(_ context.Context, input dto.CreateIngredientInput) (dto.IngredientOutput, error) {
			if input.Name != "Ingrediente Demo" || input.Unit != "KG" || !input.AverageCost.Equal(decimal.RequireFromString("4500")) {
				t.Fatalf("unexpected create input: %+v", input)
			}
			return output, nil
		},
		getFn: func(_ context.Context, id uuid.UUID) (dto.IngredientOutput, error) {
			if id != ingredientID {
				t.Fatalf("unexpected get id: %s", id)
			}
			return output, nil
		},
		listFn: func(context.Context) ([]dto.IngredientOutput, error) {
			return []dto.IngredientOutput{output}, nil
		},
		updateFn: func(_ context.Context, input dto.UpdateIngredientInput) (dto.IngredientOutput, error) {
			if input.ID != ingredientID || input.Name != "Ingrediente Editado" || !input.Stock.Equal(decimal.RequireFromString("18")) {
				t.Fatalf("unexpected update input: %+v", input)
			}
			updated := output
			updated.Name = input.Name
			updated.Stock = input.Stock
			return updated, nil
		},
		deactivateFn: func(_ context.Context, id uuid.UUID) error {
			if id != ingredientID {
				t.Fatalf("unexpected delete id: %s", id)
			}
			return nil
		},
	}

	router := ingredientTestRouter(useCase)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{
			name:       "create",
			method:     http.MethodPost,
			path:       "/ingredients/",
			body:       `{"name":"Ingrediente Demo","unit":"KG","average_cost":"4500","stock":"20","minimum_stock":"5"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "list",
			method:     http.MethodGet,
			path:       "/ingredients/",
			wantStatus: http.StatusOK,
		},
		{
			name:       "get",
			method:     http.MethodGet,
			path:       "/ingredients/" + ingredientID.String(),
			wantStatus: http.StatusOK,
		},
		{
			name:       "update",
			method:     http.MethodPut,
			path:       "/ingredients/" + ingredientID.String(),
			body:       `{"name":"Ingrediente Editado","unit":"KG","average_cost":"4800","stock":"18","minimum_stock":"4"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "delete",
			method:     http.MethodDelete,
			path:       "/ingredients/" + ingredientID.String(),
			wantStatus: http.StatusOK,
		},
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

func TestIngredientEndpointsErrors(t *testing.T) {
	ingredientID := uuid.New()
	useCase := fakeIngredientUseCase{
		createFn: func(context.Context, dto.CreateIngredientInput) (dto.IngredientOutput, error) {
			return dto.IngredientOutput{}, domainerrors.ErrInvalidInput
		},
		getFn: func(context.Context, uuid.UUID) (dto.IngredientOutput, error) {
			return dto.IngredientOutput{}, domainerrors.ErrNotFound
		},
		listFn: func(context.Context) ([]dto.IngredientOutput, error) {
			return nil, domainerrors.ErrInvalidInput
		},
		updateFn: func(context.Context, dto.UpdateIngredientInput) (dto.IngredientOutput, error) {
			return dto.IngredientOutput{}, domainerrors.ErrInactive
		},
		deactivateFn: func(context.Context, uuid.UUID) error {
			return domainerrors.ErrInactive
		},
	}

	router := ingredientTestRouter(useCase)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "create invalid decimal",
			method:     http.MethodPost,
			path:       "/ingredients/",
			body:       `{"name":"Ingrediente Demo","unit":"KG","average_cost":"not-decimal","stock":"20","minimum_stock":"5"}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_INPUT",
		},
		{
			name:       "list usecase error",
			method:     http.MethodGet,
			path:       "/ingredients/",
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_INPUT",
		},
		{
			name:       "get invalid uuid",
			method:     http.MethodGet,
			path:       "/ingredients/not-a-uuid",
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_INPUT",
		},
		{
			name:       "get not found",
			method:     http.MethodGet,
			path:       "/ingredients/" + ingredientID.String(),
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "update inactive",
			method:     http.MethodPut,
			path:       "/ingredients/" + ingredientID.String(),
			body:       `{"name":"Ingrediente Demo","unit":"KG","average_cost":"4500","stock":"20","minimum_stock":"5"}`,
			wantStatus: http.StatusConflict,
			wantCode:   "INACTIVE",
		},
		{
			name:       "delete inactive",
			method:     http.MethodDelete,
			path:       "/ingredients/" + ingredientID.String(),
			wantStatus: http.StatusConflict,
			wantCode:   "INACTIVE",
		},
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

func ingredientTestRouter(useCase fakeIngredientUseCase) http.Handler {
	return routes.NewRouter(routes.Dependencies{
		Health:      handlers.NewHealthHandler(),
		Customers:   handlers.NewCustomerHandler(nil),
		Products:    handlers.NewProductHandler(nil),
		Ingredients: handlers.NewIngredientHandler(useCase),
		Recipes:     handlers.NewRecipeHandler(nil),
		Orders:      handlers.NewOrderHandler(nil),
		Payments:    handlers.NewPaymentHandler(nil),
	})
}
