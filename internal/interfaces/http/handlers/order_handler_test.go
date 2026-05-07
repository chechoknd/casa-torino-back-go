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

type fakeOrderUseCase struct {
	createFn       func(context.Context, dto.CreateOrderInput) (dto.OrderOutput, error)
	listFn         func(context.Context, dto.ListOrdersInput) ([]dto.OrderOutput, error)
	getFn          func(context.Context, uuid.UUID) (dto.OrderOutput, error)
	addItemFn      func(context.Context, dto.AddOrderItemInput) (dto.OrderOutput, error)
	updateStatusFn func(context.Context, dto.UpdateOrderStatusInput) (dto.OrderOutput, error)
}

func (f fakeOrderUseCase) CreateOrder(ctx context.Context, input dto.CreateOrderInput) (dto.OrderOutput, error) {
	return f.createFn(ctx, input)
}

func (f fakeOrderUseCase) ListOrders(ctx context.Context, input dto.ListOrdersInput) ([]dto.OrderOutput, error) {
	return f.listFn(ctx, input)
}

func (f fakeOrderUseCase) GetOrder(ctx context.Context, id uuid.UUID) (dto.OrderOutput, error) {
	return f.getFn(ctx, id)
}

func (f fakeOrderUseCase) AddOrderItem(ctx context.Context, input dto.AddOrderItemInput) (dto.OrderOutput, error) {
	return f.addItemFn(ctx, input)
}

func (f fakeOrderUseCase) UpdateOrderStatus(ctx context.Context, input dto.UpdateOrderStatusInput) (dto.OrderOutput, error) {
	return f.updateStatusFn(ctx, input)
}

func TestOrderEndpointsSuccess(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()
	productID := uuid.New()
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	output := dto.OrderOutput{
		ID:           orderID,
		CustomerID:   customerID,
		CustomerName: "Cliente Demo",
		OrderNumber:  8,
		OrderLabel:   "#0008",
		Status:       "PENDING",
		Items: []dto.OrderItemOutput{{
			ID:          uuid.New(),
			OrderID:     orderID,
			ProductID:   productID,
			ProductName: "Producto Demo",
			Quantity:    2,
			UnitPrice:   decimal.RequireFromString("18000"),
			Subtotal:    decimal.RequireFromString("36000"),
		}},
		Subtotal:  decimal.RequireFromString("36000"),
		Discount:  decimal.RequireFromString("1000"),
		Total:     decimal.RequireFromString("35000"),
		CreatedAt: now,
		UpdatedAt: now,
	}

	useCase := fakeOrderUseCase{
		createFn: func(_ context.Context, input dto.CreateOrderInput) (dto.OrderOutput, error) {
			if input.CustomerID != customerID || !input.Discount.Equal(decimal.RequireFromString("1000")) {
				t.Fatalf("unexpected create input: %+v", input)
			}
			created := output
			created.Items = nil
			return created, nil
		},
		listFn: func(_ context.Context, input dto.ListOrdersInput) ([]dto.OrderOutput, error) {
			if input.CustomerID == nil || *input.CustomerID != customerID {
				t.Fatalf("unexpected list input: %+v", input)
			}
			return []dto.OrderOutput{output}, nil
		},
		getFn: func(_ context.Context, id uuid.UUID) (dto.OrderOutput, error) {
			if id != orderID {
				t.Fatalf("unexpected get id: %s", id)
			}
			return output, nil
		},
		addItemFn: func(_ context.Context, input dto.AddOrderItemInput) (dto.OrderOutput, error) {
			if input.OrderID != orderID || input.ProductID != productID || input.Quantity != 2 {
				t.Fatalf("unexpected add item input: %+v", input)
			}
			return output, nil
		},
		updateStatusFn: func(_ context.Context, input dto.UpdateOrderStatusInput) (dto.OrderOutput, error) {
			if input.OrderID != orderID || input.Status != "CONFIRMED" {
				t.Fatalf("unexpected update status input: %+v", input)
			}
			updated := output
			updated.Status = input.Status
			return updated, nil
		},
	}

	router := orderTestRouter(useCase)
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"create", http.MethodPost, "/orders/", `{"customer_id":"` + customerID.String() + `","discount":"1000"}`, http.StatusCreated},
		{"list", http.MethodGet, "/orders/?customer_id=" + customerID.String(), "", http.StatusOK},
		{"get", http.MethodGet, "/orders/" + orderID.String(), "", http.StatusOK},
		{"add item", http.MethodPost, "/orders/" + orderID.String() + "/items", `{"product_id":"` + productID.String() + `","quantity":2}`, http.StatusOK},
		{"update status", http.MethodPatch, "/orders/" + orderID.String() + "/status", `{"status":"CONFIRMED"}`, http.StatusOK},
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

func TestOrderEndpointsErrors(t *testing.T) {
	orderID := uuid.New()
	productID := uuid.New()
	useCase := fakeOrderUseCase{
		createFn: func(context.Context, dto.CreateOrderInput) (dto.OrderOutput, error) {
			return dto.OrderOutput{}, domainerrors.ErrInvalidInput
		},
		listFn: func(context.Context, dto.ListOrdersInput) ([]dto.OrderOutput, error) {
			return nil, domainerrors.ErrInvalidInput
		},
		getFn: func(context.Context, uuid.UUID) (dto.OrderOutput, error) {
			return dto.OrderOutput{}, domainerrors.ErrNotFound
		},
		addItemFn: func(context.Context, dto.AddOrderItemInput) (dto.OrderOutput, error) {
			return dto.OrderOutput{}, domainerrors.ErrInactive
		},
		updateStatusFn: func(context.Context, dto.UpdateOrderStatusInput) (dto.OrderOutput, error) {
			return dto.OrderOutput{}, domainerrors.ErrInvalidStatus
		},
	}

	router := orderTestRouter(useCase)
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{"create invalid customer uuid", http.MethodPost, "/orders/", `{"customer_id":"bad","discount":"1000"}`, http.StatusBadRequest, "INVALID_INPUT"},
		{"list invalid customer uuid", http.MethodGet, "/orders/?customer_id=bad", "", http.StatusBadRequest, "INVALID_INPUT"},
		{"get not found", http.MethodGet, "/orders/" + orderID.String(), "", http.StatusNotFound, "NOT_FOUND"},
		{"add item inactive", http.MethodPost, "/orders/" + orderID.String() + "/items", `{"product_id":"` + productID.String() + `","quantity":2}`, http.StatusConflict, "INACTIVE"},
		{"update invalid status", http.MethodPatch, "/orders/" + orderID.String() + "/status", `{"status":"DELIVERED"}`, http.StatusUnprocessableEntity, "INVALID_STATUS"},
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

func orderTestRouter(useCase fakeOrderUseCase) http.Handler {
	return routes.NewRouter(routes.Dependencies{
		Customers:   handlers.NewCustomerHandler(nil),
		Products:    handlers.NewProductHandler(nil),
		Ingredients: handlers.NewIngredientHandler(nil),
		Recipes:     handlers.NewRecipeHandler(nil),
		Orders:      handlers.NewOrderHandler(useCase),
		Payments:    handlers.NewPaymentHandler(nil),
	})
}
