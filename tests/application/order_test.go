package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/application/dto"
	orderuc "github.com/casatorino/backend/internal/application/usecases/order"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	"github.com/casatorino/backend/tests/mocks"
)

func TestCreateOrderSuccess(t *testing.T) {
	ct, _ := valueobjects.NewCustomerType("PERSON")
	customerRepo := &mocks.CustomerRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Customer, error) {
			return &entities.Customer{ID: uuid.New(), CustomerType: ct, IsActive: true}, nil
		},
	}
	orderRepo := &mocks.OrderRepository{CreateFn: func(context.Context, *entities.Order) error { return nil }}
	uc := orderuc.NewUseCase(orderRepo, customerRepo, &mocks.ProductRepository{})

	out, err := uc.CreateOrder(context.Background(), dto.CreateOrderInput{CustomerID: uuid.New()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Status != "PENDING" {
		t.Fatalf("unexpected status: %s", out.Status)
	}
}

func TestAddOrderItemSuccess(t *testing.T) {
	productType, _ := valueobjects.NewProductType("LUNCH")
	orderID := uuid.New()
	orderRepo := &mocks.OrderRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Order, error) {
			return &entities.Order{ID: orderID, Status: valueobjects.OrderStatusPending, Discount: decimal.Zero}, nil
		},
		AddItemFn: func(context.Context, uuid.UUID, *entities.OrderItem) error { return nil },
		UpdateFn:  func(context.Context, *entities.Order) error { return nil },
	}
	productRepo := &mocks.ProductRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Product, error) {
			return &entities.Product{ID: uuid.New(), ProductType: productType, BasePrice: decimal.RequireFromString("12000"), IsActive: true}, nil
		},
	}
	uc := orderuc.NewUseCase(orderRepo, &mocks.CustomerRepository{}, productRepo)

	out, err := uc.AddOrderItem(context.Background(), dto.AddOrderItemInput{OrderID: orderID, ProductID: uuid.New(), Quantity: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Total.Equal(decimal.RequireFromString("24000")) {
		t.Fatalf("unexpected total: %s", out.Total)
	}
}

func TestUpdateOrderStatusInvalidTransition(t *testing.T) {
	orderRepo := &mocks.OrderRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Order, error) {
			return &entities.Order{ID: uuid.New(), Status: valueobjects.OrderStatusPending}, nil
		},
	}
	uc := orderuc.NewUseCase(orderRepo, &mocks.CustomerRepository{}, &mocks.ProductRepository{})
	_, err := uc.UpdateOrderStatus(context.Background(), dto.UpdateOrderStatusInput{OrderID: uuid.New(), Status: "DELIVERED"})
	if !errors.Is(err, domainerrors.ErrInvalidStatus) {
		t.Fatalf("expected invalid status, got %v", err)
	}
}
