package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateOrderInput struct {
	CustomerID uuid.UUID
	Discount   decimal.Decimal
}

type AddOrderItemInput struct {
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Quantity  int
}

type UpdateOrderStatusInput struct {
	OrderID uuid.UUID
	Status  string
}

type OrderItemOutput struct {
	ID          uuid.UUID       `json:"id"`
	OrderID     uuid.UUID       `json:"order_id"`
	ProductID   uuid.UUID       `json:"product_id"`
	ProductName string          `json:"product_name"`
	Quantity    int             `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	Subtotal    decimal.Decimal `json:"subtotal"`
}

type OrderOutput struct {
	ID           uuid.UUID         `json:"id"`
	CustomerID   uuid.UUID         `json:"customer_id"`
	CustomerName string            `json:"customer_name"`
	OrderNumber  int64             `json:"order_number"`
	OrderLabel   string            `json:"order_label"`
	Status       string            `json:"status"`
	Items        []OrderItemOutput `json:"items"`
	Subtotal     decimal.Decimal   `json:"subtotal"`
	Discount     decimal.Decimal   `json:"discount"`
	Total        decimal.Decimal   `json:"total"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type ListOrdersInput struct {
	CustomerID *uuid.UUID
}
