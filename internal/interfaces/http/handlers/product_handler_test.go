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

type fakeProductUseCase struct {
	createFn     func(context.Context, dto.CreateProductInput) (dto.ProductOutput, error)
	getFn        func(context.Context, uuid.UUID) (dto.ProductOutput, error)
	listFn       func(context.Context, dto.ListProductsInput) ([]dto.ProductOutput, error)
	getPublicFn  func(context.Context, uuid.UUID) (dto.CatalogProductOutput, error)
	listPublicFn func(context.Context, dto.ListProductsInput) ([]dto.CatalogProductOutput, error)
	categoriesFn func(context.Context) ([]string, error)
	updateFn     func(context.Context, dto.UpdateProductInput) (dto.ProductOutput, error)
	deactivateFn func(context.Context, uuid.UUID) error
}

func (f fakeProductUseCase) CreateProduct(ctx context.Context, input dto.CreateProductInput) (dto.ProductOutput, error) {
	return f.createFn(ctx, input)
}

func (f fakeProductUseCase) GetProduct(ctx context.Context, id uuid.UUID) (dto.ProductOutput, error) {
	return f.getFn(ctx, id)
}

func (f fakeProductUseCase) ListProducts(ctx context.Context, input dto.ListProductsInput) ([]dto.ProductOutput, error) {
	return f.listFn(ctx, input)
}

func (f fakeProductUseCase) GetPublicProduct(ctx context.Context, id uuid.UUID) (dto.CatalogProductOutput, error) {
	return f.getPublicFn(ctx, id)
}

func (f fakeProductUseCase) ListPublicProducts(ctx context.Context, input dto.ListProductsInput) ([]dto.CatalogProductOutput, error) {
	return f.listPublicFn(ctx, input)
}

func (f fakeProductUseCase) ListPublicCategories(ctx context.Context) ([]string, error) {
	return f.categoriesFn(ctx)
}

func (f fakeProductUseCase) UpdateProduct(ctx context.Context, input dto.UpdateProductInput) (dto.ProductOutput, error) {
	return f.updateFn(ctx, input)
}

func (f fakeProductUseCase) DeactivateProduct(ctx context.Context, id uuid.UUID) error {
	return f.deactivateFn(ctx, id)
}

func TestProductEndpointsSuccess(t *testing.T) {
	productID := uuid.New()
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	output := dto.ProductOutput{
		ID:          productID,
		Name:        "Producto Demo",
		Description: "Producto de prueba",
		ProductType: "LUNCH",
		BasePrice:   decimal.RequireFromString("18000"),
		CostPrice:   decimal.RequireFromString("9000"),
		ImageURL:    "https://cdn.example.com/product.jpg",
		IsPublic:    true,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	useCase := fakeProductUseCase{
		createFn: func(_ context.Context, input dto.CreateProductInput) (dto.ProductOutput, error) {
			if input.Name != "Producto Demo" || input.ProductType != "LUNCH" || input.ImageURL != "https://cdn.example.com/product.jpg" || !input.IsPublic || !input.BasePrice.Equal(decimal.RequireFromString("18000")) {
				t.Fatalf("unexpected create input: %+v", input)
			}
			return output, nil
		},
		getFn: func(_ context.Context, id uuid.UUID) (dto.ProductOutput, error) {
			if id != productID {
				t.Fatalf("unexpected get id: %s", id)
			}
			return output, nil
		},
		listFn: func(_ context.Context, input dto.ListProductsInput) ([]dto.ProductOutput, error) {
			if input.ProductType != "LUNCH" {
				t.Fatalf("unexpected list input: %+v", input)
			}
			return []dto.ProductOutput{output}, nil
		},
		getPublicFn: func(_ context.Context, id uuid.UUID) (dto.CatalogProductOutput, error) {
			if id != productID {
				t.Fatalf("unexpected public get id: %s", id)
			}
			return dto.CatalogProductOutput{
				ID:          output.ID,
				Name:        output.Name,
				Description: output.Description,
				ProductType: output.ProductType,
				BasePrice:   output.BasePrice,
				ImageURL:    output.ImageURL,
				CreatedAt:   output.CreatedAt,
				UpdatedAt:   output.UpdatedAt,
			}, nil
		},
		listPublicFn: func(_ context.Context, input dto.ListProductsInput) ([]dto.CatalogProductOutput, error) {
			if input.ProductType != "LUNCH" {
				t.Fatalf("unexpected public list input: %+v", input)
			}
			return []dto.CatalogProductOutput{{
				ID:          output.ID,
				Name:        output.Name,
				Description: output.Description,
				ProductType: output.ProductType,
				BasePrice:   output.BasePrice,
				ImageURL:    output.ImageURL,
				CreatedAt:   output.CreatedAt,
				UpdatedAt:   output.UpdatedAt,
			}}, nil
		},
		categoriesFn: func(context.Context) ([]string, error) {
			return []string{"LUNCH"}, nil
		},
		updateFn: func(_ context.Context, input dto.UpdateProductInput) (dto.ProductOutput, error) {
			if input.ID != productID || input.Name != "Producto Editado" || input.ImageURL != "https://cdn.example.com/updated.jpg" || input.IsPublic || !input.BasePrice.Equal(decimal.RequireFromString("22000")) {
				t.Fatalf("unexpected update input: %+v", input)
			}
			updated := output
			updated.Name = input.Name
			updated.BasePrice = input.BasePrice
			return updated, nil
		},
		deactivateFn: func(_ context.Context, id uuid.UUID) error {
			if id != productID {
				t.Fatalf("unexpected delete id: %s", id)
			}
			return nil
		},
	}

	router := productTestRouter(useCase)

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
			path:       "/products/",
			body:       `{"name":"Producto Demo","description":"Producto de prueba","product_type":"LUNCH","base_price":"18000","cost_price":"9000","image_url":"https://cdn.example.com/product.jpg","is_public":true}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "list",
			method:     http.MethodGet,
			path:       "/products/?product_type=LUNCH",
			wantStatus: http.StatusOK,
		},
		{
			name:       "get",
			method:     http.MethodGet,
			path:       "/products/" + productID.String(),
			wantStatus: http.StatusOK,
		},
		{
			name:       "update",
			method:     http.MethodPut,
			path:       "/products/" + productID.String(),
			body:       `{"name":"Producto Editado","description":"Producto editado","product_type":"LUNCH","base_price":"22000","cost_price":"11000","image_url":"https://cdn.example.com/updated.jpg","is_public":false}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "delete",
			method:     http.MethodDelete,
			path:       "/products/" + productID.String(),
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

func TestPublicProductEndpointsSuccess(t *testing.T) {
	productID := uuid.New()
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	output := dto.CatalogProductOutput{
		ID:          productID,
		Name:        "Producto Publico",
		Description: "Producto visible",
		ProductType: "LUNCH",
		BasePrice:   decimal.RequireFromString("18000"),
		ImageURL:    "https://cdn.example.com/product.jpg",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	useCase := fakeProductUseCase{
		getPublicFn: func(_ context.Context, id uuid.UUID) (dto.CatalogProductOutput, error) {
			if id != productID {
				t.Fatalf("unexpected get id: %s", id)
			}
			return output, nil
		},
		listPublicFn: func(_ context.Context, input dto.ListProductsInput) ([]dto.CatalogProductOutput, error) {
			if input.ProductType != "LUNCH" {
				t.Fatalf("unexpected list input: %+v", input)
			}
			return []dto.CatalogProductOutput{output}, nil
		},
		categoriesFn: func(context.Context) ([]string, error) {
			return []string{"LUNCH"}, nil
		},
	}

	router := productTestRouter(useCase)

	tests := []struct {
		name string
		path string
	}{
		{name: "public list", path: "/public/products?product_type=LUNCH"},
		{name: "public get", path: "/public/products/" + productID.String()},
		{name: "public categories", path: "/public/product-categories"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := performCustomerRequest(router, http.MethodGet, tc.path, "")
			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d, body: %s", recorder.Code, http.StatusOK, recorder.Body.String())
			}
			assertSuccessEnvelope(t, recorder.Body.Bytes())
		})
	}
}

func TestProductEndpointsErrors(t *testing.T) {
	productID := uuid.New()
	useCase := fakeProductUseCase{
		createFn: func(context.Context, dto.CreateProductInput) (dto.ProductOutput, error) {
			return dto.ProductOutput{}, domainerrors.ErrInvalidInput
		},
		getFn: func(context.Context, uuid.UUID) (dto.ProductOutput, error) {
			return dto.ProductOutput{}, domainerrors.ErrNotFound
		},
		listFn: func(context.Context, dto.ListProductsInput) ([]dto.ProductOutput, error) {
			return nil, domainerrors.ErrInvalidInput
		},
		getPublicFn: func(context.Context, uuid.UUID) (dto.CatalogProductOutput, error) {
			return dto.CatalogProductOutput{}, domainerrors.ErrNotFound
		},
		listPublicFn: func(context.Context, dto.ListProductsInput) ([]dto.CatalogProductOutput, error) {
			return nil, domainerrors.ErrInvalidInput
		},
		categoriesFn: func(context.Context) ([]string, error) {
			return nil, domainerrors.ErrInvalidInput
		},
		updateFn: func(context.Context, dto.UpdateProductInput) (dto.ProductOutput, error) {
			return dto.ProductOutput{}, domainerrors.ErrInactive
		},
		deactivateFn: func(context.Context, uuid.UUID) error {
			return domainerrors.ErrInactive
		},
	}

	router := productTestRouter(useCase)

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
			path:       "/products/",
			body:       `{"name":"Producto Demo","description":"Producto de prueba","product_type":"LUNCH","base_price":"not-decimal","cost_price":"9000"}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_INPUT",
		},
		{
			name:       "list usecase error",
			method:     http.MethodGet,
			path:       "/products/",
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_INPUT",
		},
		{
			name:       "get invalid uuid",
			method:     http.MethodGet,
			path:       "/products/not-a-uuid",
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_INPUT",
		},
		{
			name:       "get not found",
			method:     http.MethodGet,
			path:       "/products/" + productID.String(),
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "update inactive",
			method:     http.MethodPut,
			path:       "/products/" + productID.String(),
			body:       `{"name":"Producto Demo","description":"Producto de prueba","product_type":"LUNCH","base_price":"18000","cost_price":"9000"}`,
			wantStatus: http.StatusConflict,
			wantCode:   "INACTIVE",
		},
		{
			name:       "delete inactive",
			method:     http.MethodDelete,
			path:       "/products/" + productID.String(),
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

func productTestRouter(useCase fakeProductUseCase) http.Handler {
	return routes.NewRouter(routes.Dependencies{
		Customers:   handlers.NewCustomerHandler(nil),
		Products:    handlers.NewProductHandler(useCase),
		Ingredients: handlers.NewIngredientHandler(nil),
		Recipes:     handlers.NewRecipeHandler(nil),
		Orders:      handlers.NewOrderHandler(nil),
		Payments:    handlers.NewPaymentHandler(nil),
	})
}
