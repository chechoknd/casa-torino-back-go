package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/application/dto"
	paymentuc "github.com/casatorino/backend/internal/application/usecases/payment"
	"github.com/casatorino/backend/internal/domain/entities"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	"github.com/casatorino/backend/tests/mocks"
)

func TestCreatePaymentSuccess(t *testing.T) {
	orderRepo := &mocks.OrderRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Order, error) {
			return &entities.Order{ID: uuid.New(), Status: valueobjects.OrderStatusPending}, nil
		},
	}
	paymentRepo := &mocks.PaymentRepository{
		CreateFn: func(context.Context, *entities.Payment) error { return nil },
	}
	uc := paymentuc.NewUseCase(paymentRepo, orderRepo)

	out, err := uc.CreatePayment(context.Background(), dto.CreatePaymentInput{
		OrderID: uuid.New(), Amount: decimal.RequireFromString("10000"), Method: "CASH", Status: "PENDING",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Method != "CASH" {
		t.Fatalf("unexpected method: %s", out.Method)
	}
}

func TestGetPaymentsByOrderSuccess(t *testing.T) {
	orderID := uuid.New()
	orderRepo := &mocks.OrderRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Order, error) {
			return &entities.Order{ID: orderID, Status: valueobjects.OrderStatusPending}, nil
		},
	}
	method, _ := valueobjects.NewPaymentMethod("CASH")
	status, _ := valueobjects.NewPaymentStatus("PAID")
	paymentRepo := &mocks.PaymentRepository{
		ListByOrderIDFn: func(context.Context, uuid.UUID) ([]entities.Payment, error) {
			return []entities.Payment{{ID: uuid.New(), OrderID: orderID, Method: method, Status: status, Amount: decimal.RequireFromString("1000")}}, nil
		},
	}
	uc := paymentuc.NewUseCase(paymentRepo, orderRepo)
	items, err := uc.GetPaymentsByOrder(context.Background(), orderID)
	if err != nil || len(items) != 1 {
		t.Fatalf("unexpected result: %v len=%d", err, len(items))
	}
}

func TestUpdatePaymentStatusSuccess(t *testing.T) {
	method, _ := valueobjects.NewPaymentMethod("CASH")
	status, _ := valueobjects.NewPaymentStatus("PENDING")
	paymentRepo := &mocks.PaymentRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Payment, error) {
			return &entities.Payment{ID: uuid.New(), Method: method, Status: status}, nil
		},
		UpdateFn: func(context.Context, *entities.Payment) error { return nil },
	}
	uc := paymentuc.NewUseCase(paymentRepo, &mocks.OrderRepository{})
	out, err := uc.UpdatePaymentStatus(context.Background(), dto.UpdatePaymentStatusInput{PaymentID: uuid.New(), Status: "PAID"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != "PAID" {
		t.Fatalf("unexpected status: %s", out.Status)
	}
}
