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
	orderID := uuid.New()
	productID := uuid.New()
	orderRepo := &mocks.OrderRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Order, error) {
			return &entities.Order{
				ID:          orderID,
				OrderNumber: 7,
				Status:      valueobjects.OrderStatusPending,
				Items: []entities.OrderItem{
					{ProductID: productID, Quantity: 2},
				},
			}, nil
		},
	}
	paymentRepo := &mocks.PaymentRepository{
		CreateFn: func(context.Context, *entities.Payment) error { return nil },
	}
	productType, _ := valueobjects.NewProductType("LUNCH")
	productRepo := &mocks.ProductRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Product, error) {
			return &entities.Product{ID: productID, Name: "Menu del dia", ProductType: productType, IsActive: true}, nil
		},
	}
	uc := paymentuc.NewUseCase(paymentRepo, orderRepo, productRepo)

	out, err := uc.CreatePayment(context.Background(), dto.CreatePaymentInput{
		OrderID: orderID, Amount: decimal.RequireFromString("10000"), Method: "CASH", Status: "PENDING",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Method != "CASH" {
		t.Fatalf("unexpected method: %s", out.Method)
	}
	if out.OrderLabel != "#0007" || len(out.Products) != 1 || out.Products[0].ProductName != "Menu del dia" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestGetPaymentsByOrderSuccess(t *testing.T) {
	orderID := uuid.New()
	productID := uuid.New()
	orderRepo := &mocks.OrderRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Order, error) {
			return &entities.Order{
				ID:          orderID,
				OrderNumber: 3,
				Status:      valueobjects.OrderStatusPending,
				Items:       []entities.OrderItem{{ProductID: productID, Quantity: 1}},
			}, nil
		},
	}
	method, _ := valueobjects.NewPaymentMethod("CASH")
	status, _ := valueobjects.NewPaymentStatus("PAID")
	paymentRepo := &mocks.PaymentRepository{
		ListByOrderIDFn: func(context.Context, uuid.UUID) ([]entities.Payment, error) {
			return []entities.Payment{{ID: uuid.New(), OrderID: orderID, Method: method, Status: status, Amount: decimal.RequireFromString("1000")}}, nil
		},
	}
	productType, _ := valueobjects.NewProductType("LUNCH")
	productRepo := &mocks.ProductRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Product, error) {
			return &entities.Product{ID: productID, Name: "Jugo natural", ProductType: productType, IsActive: true}, nil
		},
	}
	uc := paymentuc.NewUseCase(paymentRepo, orderRepo, productRepo)
	items, err := uc.GetPaymentsByOrder(context.Background(), orderID)
	if err != nil || len(items) != 1 {
		t.Fatalf("unexpected result: %v len=%d", err, len(items))
	}
	if items[0].Products[0].ProductName != "Jugo natural" {
		t.Fatalf("unexpected products: %+v", items[0].Products)
	}
}

func TestUpdatePaymentStatusSuccess(t *testing.T) {
	orderID := uuid.New()
	productID := uuid.New()
	method, _ := valueobjects.NewPaymentMethod("CASH")
	status, _ := valueobjects.NewPaymentStatus("PENDING")
	paymentRepo := &mocks.PaymentRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Payment, error) {
			return &entities.Payment{ID: uuid.New(), OrderID: orderID, Method: method, Status: status}, nil
		},
		UpdateFn: func(context.Context, *entities.Payment) error { return nil },
	}
	orderRepo := &mocks.OrderRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Order, error) {
			return &entities.Order{
				ID:          orderID,
				OrderNumber: 9,
				Status:      valueobjects.OrderStatusPending,
				Items:       []entities.OrderItem{{ProductID: productID, Quantity: 1}},
			}, nil
		},
	}
	productType, _ := valueobjects.NewProductType("LUNCH")
	productRepo := &mocks.ProductRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Product, error) {
			return &entities.Product{ID: productID, Name: "Torta zanahoria", ProductType: productType, IsActive: true}, nil
		},
	}
	uc := paymentuc.NewUseCase(paymentRepo, orderRepo, productRepo)
	out, err := uc.UpdatePaymentStatus(context.Background(), dto.UpdatePaymentStatusInput{PaymentID: uuid.New(), Status: "PAID"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != "PAID" {
		t.Fatalf("unexpected status: %s", out.Status)
	}
}
