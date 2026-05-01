package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	customeruc "github.com/casatorino/backend/internal/application/usecases/customer"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	"github.com/casatorino/backend/tests/mocks"
)

func TestCreateCustomerSuccess(t *testing.T) {
	repo := &mocks.CustomerRepository{
		FindByEmailFn: func(context.Context, string) (*entities.Customer, error) { return nil, domainerrors.ErrNotFound },
		CreateFn:      func(context.Context, *entities.Customer) error { return nil },
	}
	uc := customeruc.NewUseCase(repo)

	out, err := uc.CreateCustomer(context.Background(), dto.CreateCustomerInput{
		FullName:     "Ana",
		Phone:        "300",
		Email:        "ana@example.com",
		CustomerType: "PERSON",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Email != "ana@example.com" {
		t.Fatalf("unexpected email: %s", out.Email)
	}
}

func TestCreateCustomerDuplicateEmail(t *testing.T) {
	repo := &mocks.CustomerRepository{
		FindByEmailFn: func(context.Context, string) (*entities.Customer, error) {
			return &entities.Customer{ID: uuid.New(), IsActive: true}, nil
		},
		CreateFn: func(context.Context, *entities.Customer) error { return nil },
	}
	uc := customeruc.NewUseCase(repo)

	_, err := uc.CreateCustomer(context.Background(), dto.CreateCustomerInput{
		FullName: "Ana", Phone: "300", Email: "ana@example.com", CustomerType: "PERSON",
	})
	if !errors.Is(err, domainerrors.ErrDuplicateEmail) {
		t.Fatalf("expected duplicate email, got %v", err)
	}
}

func TestGetCustomerInactive(t *testing.T) {
	repo := &mocks.CustomerRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Customer, error) {
			return &entities.Customer{ID: uuid.New(), IsActive: false}, nil
		},
	}
	uc := customeruc.NewUseCase(repo)
	_, err := uc.GetCustomer(context.Background(), uuid.New())
	if !errors.Is(err, domainerrors.ErrInactive) {
		t.Fatalf("expected inactive error, got %v", err)
	}
}

func TestUpdateCustomerSuccess(t *testing.T) {
	id := uuid.New()
	repo := &mocks.CustomerRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Customer, error) {
			ct, _ := valueobjects.NewCustomerType("PERSON")
			return &entities.Customer{ID: id, IsActive: true, CustomerType: ct, Email: "old@example.com"}, nil
		},
		FindByEmailFn: func(context.Context, string) (*entities.Customer, error) { return nil, domainerrors.ErrNotFound },
		UpdateFn:      func(context.Context, *entities.Customer) error { return nil },
	}
	uc := customeruc.NewUseCase(repo)
	out, err := uc.UpdateCustomer(context.Background(), dto.UpdateCustomerInput{
		ID: id, FullName: "Nuevo", Phone: "123", Email: "new@example.com", CustomerType: "COMPANY",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.CustomerType != "COMPANY" {
		t.Fatalf("unexpected type: %s", out.CustomerType)
	}
}

func TestDeactivateCustomerSuccess(t *testing.T) {
	called := false
	repo := &mocks.CustomerRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Customer, error) {
			return &entities.Customer{ID: uuid.New(), IsActive: true}, nil
		},
		DeactivateFn: func(context.Context, uuid.UUID, time.Time) error {
			called = true
			return nil
		},
	}
	uc := customeruc.NewUseCase(repo)
	if err := uc.DeactivateCustomer(context.Background(), uuid.New()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected deactivate call")
	}
}
