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

type fakePaymentUseCase struct {
	createFn       func(context.Context, dto.CreatePaymentInput) (dto.PaymentOutput, error)
	getByOrderFn   func(context.Context, uuid.UUID) ([]dto.PaymentOutput, error)
	listFn         func(context.Context) ([]dto.PaymentOutput, error)
	updateStatusFn func(context.Context, dto.UpdatePaymentStatusInput) (dto.PaymentOutput, error)
}

func (f fakePaymentUseCase) CreatePayment(ctx context.Context, input dto.CreatePaymentInput) (dto.PaymentOutput, error) {
	return f.createFn(ctx, input)
}

func (f fakePaymentUseCase) GetPaymentsByOrder(ctx context.Context, orderID uuid.UUID) ([]dto.PaymentOutput, error) {
	return f.getByOrderFn(ctx, orderID)
}

func (f fakePaymentUseCase) ListPayments(ctx context.Context) ([]dto.PaymentOutput, error) {
	return f.listFn(ctx)
}

func (f fakePaymentUseCase) UpdatePaymentStatus(ctx context.Context, input dto.UpdatePaymentStatusInput) (dto.PaymentOutput, error) {
	return f.updateStatusFn(ctx, input)
}

func TestPaymentEndpointsSuccess(t *testing.T) {
	paymentID := uuid.New()
	orderID := uuid.New()
	productID := uuid.New()
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	output := dto.PaymentOutput{
		ID:          paymentID,
		OrderID:     orderID,
		OrderNumber: 8,
		OrderLabel:  "#0008",
		Amount:      decimal.RequireFromString("35000"),
		Method:      "CASH",
		Status:      "PENDING",
		Products: []dto.PaymentProductOutput{{
			ProductID:   productID,
			ProductName: "Producto Demo",
			Quantity:    2,
		}},
		CreatedAt: now,
		UpdatedAt: now,
	}

	useCase := fakePaymentUseCase{
		createFn: func(_ context.Context, input dto.CreatePaymentInput) (dto.PaymentOutput, error) {
			if input.OrderID != orderID || input.Method != "CASH" || input.Status != "PENDING" || !input.Amount.Equal(decimal.RequireFromString("35000")) {
				t.Fatalf("unexpected create input: %+v", input)
			}
			return output, nil
		},
		getByOrderFn: func(_ context.Context, id uuid.UUID) ([]dto.PaymentOutput, error) {
			if id != orderID {
				t.Fatalf("unexpected order id: %s", id)
			}
			return []dto.PaymentOutput{output}, nil
		},
		listFn: func(context.Context) ([]dto.PaymentOutput, error) {
			return []dto.PaymentOutput{output}, nil
		},
		updateStatusFn: func(_ context.Context, input dto.UpdatePaymentStatusInput) (dto.PaymentOutput, error) {
			if input.PaymentID != paymentID || input.Status != "PAID" {
				t.Fatalf("unexpected update status input: %+v", input)
			}
			updated := output
			updated.Status = input.Status
			return updated, nil
		},
	}

	router := paymentTestRouter(useCase)
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"list", http.MethodGet, "/payments/", "", http.StatusOK},
		{"create", http.MethodPost, "/payments/", `{"order_id":"` + orderID.String() + `","amount":"35000","method":"CASH","status":"PENDING"}`, http.StatusCreated},
		{"update status", http.MethodPatch, "/payments/" + paymentID.String() + "/status", `{"status":"PAID"}`, http.StatusOK},
		{"get by order", http.MethodGet, "/orders/" + orderID.String() + "/payments", "", http.StatusOK},
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

func TestPaymentEndpointsErrors(t *testing.T) {
	paymentID := uuid.New()
	orderID := uuid.New()
	useCase := fakePaymentUseCase{
		createFn: func(context.Context, dto.CreatePaymentInput) (dto.PaymentOutput, error) {
			return dto.PaymentOutput{}, domainerrors.ErrInvalidInput
		},
		getByOrderFn: func(context.Context, uuid.UUID) ([]dto.PaymentOutput, error) {
			return nil, domainerrors.ErrNotFound
		},
		listFn: func(context.Context) ([]dto.PaymentOutput, error) {
			return nil, domainerrors.ErrInvalidInput
		},
		updateStatusFn: func(context.Context, dto.UpdatePaymentStatusInput) (dto.PaymentOutput, error) {
			return dto.PaymentOutput{}, domainerrors.ErrInvalidInput
		},
	}

	router := paymentTestRouter(useCase)
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{"list usecase error", http.MethodGet, "/payments/", "", http.StatusBadRequest, "INVALID_INPUT"},
		{"create invalid order uuid", http.MethodPost, "/payments/", `{"order_id":"bad","amount":"35000","method":"CASH","status":"PENDING"}`, http.StatusBadRequest, "INVALID_INPUT"},
		{"create invalid amount", http.MethodPost, "/payments/", `{"order_id":"` + orderID.String() + `","amount":"bad","method":"CASH","status":"PENDING"}`, http.StatusBadRequest, "INVALID_INPUT"},
		{"update invalid uuid", http.MethodPatch, "/payments/not-a-uuid/status", `{"status":"PAID"}`, http.StatusBadRequest, "INVALID_INPUT"},
		{"update usecase error", http.MethodPatch, "/payments/" + paymentID.String() + "/status", `{"status":"BAD"}`, http.StatusBadRequest, "INVALID_INPUT"},
		{"get by order not found", http.MethodGet, "/orders/" + orderID.String() + "/payments", "", http.StatusNotFound, "NOT_FOUND"},
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

func paymentTestRouter(useCase fakePaymentUseCase) http.Handler {
	return routes.NewRouter(routes.Dependencies{
		Customers:   handlers.NewCustomerHandler(nil),
		Products:    handlers.NewProductHandler(nil),
		Ingredients: handlers.NewIngredientHandler(nil),
		Recipes:     handlers.NewRecipeHandler(nil),
		Orders:      handlers.NewOrderHandler(nil),
		Payments:    handlers.NewPaymentHandler(useCase),
	})
}
