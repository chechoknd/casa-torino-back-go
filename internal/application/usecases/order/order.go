package order

import (
	"context"
	"fmt"
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
	orders    repositories.OrderRepository
	customers repositories.CustomerRepository
	products  repositories.ProductRepository
}

func NewUseCase(orders repositories.OrderRepository, customers repositories.CustomerRepository, products repositories.ProductRepository) *UseCase {
	return &UseCase{
		orders:    orders,
		customers: customers,
		products:  products,
	}
}

func (uc *UseCase) CreateOrder(ctx context.Context, input dto.CreateOrderInput) (dto.OrderOutput, error) {
	customer, err := uc.customers.FindByID(ctx, input.CustomerID)
	if err != nil {
		return dto.OrderOutput{}, err
	}
	if !customer.IsActive {
		return dto.OrderOutput{}, domainerrors.ErrInactive
	}
	if input.Discount.IsNegative() {
		return dto.OrderOutput{}, domainerrors.ErrInvalidInput
	}

	now := time.Now().UTC()
	order := &entities.Order{
		ID:          uuid.New(),
		CustomerID:  input.CustomerID,
		OrderNumber: 0,
		Status:      valueobjects.OrderStatusPending,
		Items:       []entities.OrderItem{},
		Subtotal:    decimal.Zero,
		Discount:    input.Discount,
		Total:       decimal.Zero,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.orders.Create(ctx, order); err != nil {
		return dto.OrderOutput{}, err
	}

	return uc.toOrderOutput(ctx, *order)
}

func (uc *UseCase) AddOrderItem(ctx context.Context, input dto.AddOrderItemInput) (dto.OrderOutput, error) {
	if input.Quantity <= 0 {
		return dto.OrderOutput{}, domainerrors.ErrInvalidInput
	}

	order, err := uc.orders.FindByID(ctx, input.OrderID)
	if err != nil {
		return dto.OrderOutput{}, err
	}

	product, err := uc.products.FindByID(ctx, input.ProductID)
	if err != nil {
		return dto.OrderOutput{}, err
	}
	if !product.IsActive {
		return dto.OrderOutput{}, domainerrors.ErrInactive
	}

	item := &entities.OrderItem{
		ID:        uuid.New(),
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  input.Quantity,
		UnitPrice: product.BasePrice,
		Subtotal:  product.BasePrice.Mul(decimal.NewFromInt(int64(input.Quantity))),
	}

	if err := uc.orders.AddItem(ctx, order.ID, item); err != nil {
		return dto.OrderOutput{}, err
	}

	order.Items = append(order.Items, *item)
	order.CalculateTotal()
	order.UpdatedAt = time.Now().UTC()

	if err := uc.orders.Update(ctx, order); err != nil {
		return dto.OrderOutput{}, err
	}

	return uc.toOrderOutput(ctx, *order)
}

func (uc *UseCase) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (dto.OrderOutput, error) {
	order, err := uc.orders.FindByID(ctx, orderID)
	if err != nil {
		return dto.OrderOutput{}, err
	}

	order.CalculateTotal()
	order.UpdatedAt = time.Now().UTC()

	if err := uc.orders.Update(ctx, order); err != nil {
		return dto.OrderOutput{}, err
	}

	return uc.toOrderOutput(ctx, *order)
}

func (uc *UseCase) UpdateOrderStatus(ctx context.Context, input dto.UpdateOrderStatusInput) (dto.OrderOutput, error) {
	order, err := uc.orders.FindByID(ctx, input.OrderID)
	if err != nil {
		return dto.OrderOutput{}, err
	}

	status, err := valueobjects.NewOrderStatus(input.Status)
	if err != nil {
		return dto.OrderOutput{}, err
	}

	if !order.CanTransitionTo(status) {
		return dto.OrderOutput{}, domainerrors.ErrInvalidStatus
	}

	order.Status = status
	order.UpdatedAt = time.Now().UTC()

	if err := uc.orders.Update(ctx, order); err != nil {
		return dto.OrderOutput{}, err
	}

	return uc.toOrderOutput(ctx, *order)
}

func (uc *UseCase) GetOrder(ctx context.Context, id uuid.UUID) (dto.OrderOutput, error) {
	order, err := uc.orders.FindByID(ctx, id)
	if err != nil {
		return dto.OrderOutput{}, err
	}
	return uc.toOrderOutput(ctx, *order)
}

func (uc *UseCase) ListOrders(ctx context.Context, input dto.ListOrdersInput) ([]dto.OrderOutput, error) {
	var (
		orders []entities.Order
		err    error
	)

	if input.CustomerID != nil {
		orders, err = uc.orders.ListByCustomerID(ctx, *input.CustomerID)
	} else {
		orders, err = uc.orders.List(ctx)
	}
	if err != nil {
		return nil, err
	}

	output := make([]dto.OrderOutput, 0, len(orders))
	for _, order := range orders {
		mapped, err := uc.toOrderOutput(ctx, order)
		if err != nil {
			return nil, err
		}
		output = append(output, mapped)
	}

	return output, nil
}

func (uc *UseCase) toOrderOutput(ctx context.Context, order entities.Order) (dto.OrderOutput, error) {
	customerName := ""
	customer, err := uc.customers.FindByID(ctx, order.CustomerID)
	if err == nil {
		customerName = customer.FullName
	}

	items := make([]dto.OrderItemOutput, 0, len(order.Items))
	for _, item := range order.Items {
		productName := ""
		product, err := uc.products.FindByID(ctx, item.ProductID)
		if err == nil {
			productName = product.Name
		}

		items = append(items, dto.OrderItemOutput{
			ID:          item.ID,
			OrderID:     item.OrderID,
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
		})
	}

	return dto.OrderOutput{
		ID:           order.ID,
		CustomerID:   order.CustomerID,
		CustomerName: customerName,
		OrderNumber:  order.OrderNumber,
		OrderLabel:   formatOrderLabel(order.OrderNumber),
		Status:       string(order.Status),
		Items:        items,
		Subtotal:     order.Subtotal,
		Discount:     order.Discount,
		Total:        order.Total,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}, nil
}

func formatOrderLabel(orderNumber int64) string {
	if orderNumber <= 0 {
		return ""
	}

	return fmt.Sprintf("#%04d", orderNumber)
}
