package customerpanel

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	"github.com/casatorino/backend/tests/mocks"
)

type fakeOrderLister struct {
	listFn func(context.Context, dto.ListOrdersInput) ([]dto.OrderOutput, error)
}

func (f fakeOrderLister) ListOrders(ctx context.Context, input dto.ListOrdersInput) ([]dto.OrderOutput, error) {
	return f.listFn(ctx, input)
}

func TestGetProfileFindsCustomerByAuthenticatedEmail(t *testing.T) {
	customerID := uuid.New()
	now := time.Date(2026, 5, 13, 12, 0, 0, 0, time.UTC)
	repo := &mocks.CustomerRepository{
		FindByEmailFn: func(_ context.Context, email string) (*entities.Customer, error) {
			if email != "laura@example.com" {
				t.Fatalf("email = %q", email)
			}
			return &entities.Customer{
				ID:           customerID,
				FullName:     "Laura",
				Phone:        "3001001001",
				Email:        email,
				CustomerType: valueobjects.CustomerTypePerson,
				CreatedAt:    now,
				UpdatedAt:    now,
				IsActive:     true,
			}, nil
		},
	}
	uc := NewUseCase(repo, fakeOrderLister{})

	out, err := uc.GetProfile(context.Background(), " Laura@Example.COM ")
	if err != nil {
		t.Fatalf("GetProfile error = %v", err)
	}
	if out.ID != customerID || out.Email != "laura@example.com" {
		t.Fatalf("unexpected profile: %+v", out)
	}
}

func TestListOrdersUsesResolvedCustomerID(t *testing.T) {
	customerID := uuid.New()
	repo := &mocks.CustomerRepository{
		FindByEmailFn: func(context.Context, string) (*entities.Customer, error) {
			return &entities.Customer{
				ID:           customerID,
				Email:        "laura@example.com",
				CustomerType: valueobjects.CustomerTypePerson,
				IsActive:     true,
			}, nil
		},
	}
	uc := NewUseCase(repo, fakeOrderLister{
		listFn: func(_ context.Context, input dto.ListOrdersInput) ([]dto.OrderOutput, error) {
			if input.CustomerID == nil || *input.CustomerID != customerID {
				t.Fatalf("unexpected order filter: %+v", input)
			}
			return []dto.OrderOutput{{ID: uuid.New(), CustomerID: customerID}}, nil
		},
	})

	orders, err := uc.ListOrders(context.Background(), "laura@example.com")
	if err != nil {
		t.Fatalf("ListOrders error = %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("orders len = %d", len(orders))
	}
}

func TestGetProfileRejectsMissingCustomer(t *testing.T) {
	repo := &mocks.CustomerRepository{
		FindByEmailFn: func(context.Context, string) (*entities.Customer, error) {
			return nil, domainerrors.ErrNotFound
		},
	}
	uc := NewUseCase(repo, fakeOrderLister{})

	_, err := uc.GetProfile(context.Background(), "missing@example.com")
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Fatalf("error = %v, want ErrNotFound", err)
	}
}
