package customer

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/repositories"
	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type UseCase struct {
	customers repositories.CustomerRepository
}

func NewUseCase(customers repositories.CustomerRepository) *UseCase {
	return &UseCase{customers: customers}
}

func (uc *UseCase) CreateCustomer(ctx context.Context, input dto.CreateCustomerInput) (dto.CustomerOutput, error) {
	customerType, err := valueobjects.NewCustomerType(input.CustomerType)
	if err != nil {
		return dto.CustomerOutput{}, err
	}

	if strings.TrimSpace(input.FullName) == "" || strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Phone) == "" {
		return dto.CustomerOutput{}, domainerrors.ErrInvalidInput
	}

	existing, err := uc.customers.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return dto.CustomerOutput{}, err
	}
	if existing != nil {
		return dto.CustomerOutput{}, domainerrors.ErrDuplicateEmail
	}

	now := time.Now().UTC()
	customer := &entities.Customer{
		ID:           uuid.New(),
		FullName:     strings.TrimSpace(input.FullName),
		Phone:        strings.TrimSpace(input.Phone),
		Email:        strings.TrimSpace(strings.ToLower(input.Email)),
		CustomerType: customerType,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}

	if err := uc.customers.Create(ctx, customer); err != nil {
		return dto.CustomerOutput{}, err
	}

	return toCustomerOutput(*customer), nil
}

func (uc *UseCase) GetCustomer(ctx context.Context, id uuid.UUID) (dto.CustomerOutput, error) {
	customer, err := uc.customers.FindByID(ctx, id)
	if err != nil {
		return dto.CustomerOutput{}, err
	}
	if customer == nil || !customer.IsActive {
		return dto.CustomerOutput{}, domainerrors.ErrInactive
	}

	return toCustomerOutput(*customer), nil
}

func (uc *UseCase) ListCustomers(ctx context.Context) ([]dto.CustomerOutput, error) {
	customers, err := uc.customers.List(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]dto.CustomerOutput, 0, len(customers))
	for _, customer := range customers {
		if !customer.IsActive {
			continue
		}
		output = append(output, toCustomerOutput(customer))
	}

	return output, nil
}

func (uc *UseCase) UpdateCustomer(ctx context.Context, input dto.UpdateCustomerInput) (dto.CustomerOutput, error) {
	customer, err := uc.customers.FindByID(ctx, input.ID)
	if err != nil {
		return dto.CustomerOutput{}, err
	}
	if !customer.IsActive {
		return dto.CustomerOutput{}, domainerrors.ErrInactive
	}

	customerType, err := valueobjects.NewCustomerType(input.CustomerType)
	if err != nil {
		return dto.CustomerOutput{}, err
	}

	existing, err := uc.customers.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return dto.CustomerOutput{}, err
	}
	if existing != nil && existing.ID != customer.ID {
		return dto.CustomerOutput{}, domainerrors.ErrDuplicateEmail
	}

	if strings.TrimSpace(input.FullName) == "" || strings.TrimSpace(input.Phone) == "" || strings.TrimSpace(input.Email) == "" {
		return dto.CustomerOutput{}, domainerrors.ErrInvalidInput
	}

	customer.FullName = strings.TrimSpace(input.FullName)
	customer.Phone = strings.TrimSpace(input.Phone)
	customer.Email = strings.TrimSpace(strings.ToLower(input.Email))
	customer.CustomerType = customerType
	customer.UpdatedAt = time.Now().UTC()

	if err := uc.customers.Update(ctx, customer); err != nil {
		return dto.CustomerOutput{}, err
	}

	return toCustomerOutput(*customer), nil
}

func (uc *UseCase) DeactivateCustomer(ctx context.Context, id uuid.UUID) error {
	customer, err := uc.customers.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if customer == nil || !customer.IsActive {
		return domainerrors.ErrInactive
	}

	return uc.customers.Deactivate(ctx, id, time.Now().UTC())
}

func toCustomerOutput(customer entities.Customer) dto.CustomerOutput {
	return dto.CustomerOutput{
		ID:           customer.ID,
		FullName:     customer.FullName,
		Phone:        customer.Phone,
		Email:        customer.Email,
		CustomerType: string(customer.CustomerType),
		CreatedAt:    customer.CreatedAt,
		UpdatedAt:    customer.UpdatedAt,
		IsActive:     customer.IsActive,
	}
}
