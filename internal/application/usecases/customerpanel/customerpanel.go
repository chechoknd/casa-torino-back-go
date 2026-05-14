package customerpanel

import (
	"context"
	"strings"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/repositories"
)

type OrderLister interface {
	ListOrders(ctx context.Context, input dto.ListOrdersInput) ([]dto.OrderOutput, error)
}

type UseCase struct {
	customers repositories.CustomerRepository
	orders    OrderLister
}

func NewUseCase(customers repositories.CustomerRepository, orders OrderLister) *UseCase {
	return &UseCase{customers: customers, orders: orders}
}

func (uc *UseCase) GetProfile(ctx context.Context, userEmail string) (dto.CustomerOutput, error) {
	customer, err := uc.customerByEmail(ctx, userEmail)
	if err != nil {
		return dto.CustomerOutput{}, err
	}

	return dto.CustomerOutput{
		ID:           customer.ID,
		FullName:     customer.FullName,
		Phone:        customer.Phone,
		Email:        customer.Email,
		CustomerType: string(customer.CustomerType),
		CreatedAt:    customer.CreatedAt,
		UpdatedAt:    customer.UpdatedAt,
		IsActive:     customer.IsActive,
	}, nil
}

func (uc *UseCase) ListOrders(ctx context.Context, userEmail string) ([]dto.OrderOutput, error) {
	customer, err := uc.customerByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return uc.orders.ListOrders(ctx, dto.ListOrdersInput{CustomerID: &customer.ID})
}

func (uc *UseCase) customerByEmail(ctx context.Context, userEmail string) (*entities.Customer, error) {
	email := strings.TrimSpace(strings.ToLower(userEmail))
	if email == "" {
		return nil, domainerrors.ErrUnauthorized
	}

	customer, err := uc.customers.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if customer == nil || !customer.IsActive {
		return nil, domainerrors.ErrNotFound
	}

	return customer, nil
}
