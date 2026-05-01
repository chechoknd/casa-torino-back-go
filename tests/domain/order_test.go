package domain_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/domain/entities"
	"github.com/casatorino/backend/internal/domain/valueobjects"
)

func TestOrderCanTransitionTo(t *testing.T) {
	t.Parallel()

	order := entities.Order{Status: valueobjects.OrderStatusPending}

	if !order.CanTransitionTo(valueobjects.OrderStatusConfirmed) {
		t.Fatalf("expected pending to transition to confirmed")
	}

	if order.CanTransitionTo(valueobjects.OrderStatusDelivered) {
		t.Fatalf("expected pending to reject delivered")
	}
}

func TestOrderCalculateTotal(t *testing.T) {
	t.Parallel()

	order := entities.Order{
		Discount: decimal.RequireFromString("500"),
		Items: []entities.OrderItem{
			{Quantity: 2, UnitPrice: decimal.RequireFromString("12000")},
			{Quantity: 1, UnitPrice: decimal.RequireFromString("8000")},
		},
	}

	order.CalculateTotal()

	if !order.Subtotal.Equal(decimal.RequireFromString("32000")) {
		t.Fatalf("unexpected subtotal: %s", order.Subtotal)
	}

	if !order.Total.Equal(decimal.RequireFromString("31500")) {
		t.Fatalf("unexpected total: %s", order.Total)
	}

	if !order.Items[0].Subtotal.Equal(decimal.RequireFromString("24000")) {
		t.Fatalf("unexpected first item subtotal: %s", order.Items[0].Subtotal)
	}
}

func TestOrderCalculateTotalDoesNotGoNegative(t *testing.T) {
	t.Parallel()

	order := entities.Order{
		Discount: decimal.RequireFromString("9999"),
		Items: []entities.OrderItem{
			{Quantity: 1, UnitPrice: decimal.RequireFromString("5000")},
		},
	}

	order.CalculateTotal()

	if !order.Total.Equal(decimal.Zero) {
		t.Fatalf("expected total to clamp at zero, got %s", order.Total)
	}
}
