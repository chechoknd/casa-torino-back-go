package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type Order struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	OrderNumber int64
	Status      valueobjects.OrderStatus
	Items       []OrderItem
	Subtotal    decimal.Decimal
	Discount    decimal.Decimal
	Total       decimal.Decimal
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (o *Order) CalculateTotal() {
	subtotal := decimal.Zero

	for index := range o.Items {
		itemSubtotal := o.Items[index].UnitPrice.Mul(decimal.NewFromInt(int64(o.Items[index].Quantity)))
		o.Items[index].Subtotal = itemSubtotal
		subtotal = subtotal.Add(itemSubtotal)
	}

	o.Subtotal = subtotal
	o.Total = subtotal.Sub(o.Discount)
	if o.Total.IsNegative() {
		o.Total = decimal.Zero
	}
}

func (o Order) CanTransitionTo(status valueobjects.OrderStatus) bool {
	return o.Status.CanTransitionTo(status)
}
