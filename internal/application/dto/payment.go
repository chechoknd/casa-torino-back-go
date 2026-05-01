package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreatePaymentInput struct {
	OrderID uuid.UUID
	Amount  decimal.Decimal
	Method  string
	Status  string
}

type UpdatePaymentStatusInput struct {
	PaymentID uuid.UUID
	Status    string
}

type PaymentProductOutput struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
}

type PaymentOutput struct {
	ID          uuid.UUID              `json:"id"`
	OrderID     uuid.UUID              `json:"order_id"`
	OrderNumber int64                  `json:"order_number"`
	OrderLabel  string                 `json:"order_label"`
	Amount      decimal.Decimal        `json:"amount"`
	Method      string                 `json:"method"`
	Status      string                 `json:"status"`
	Products    []PaymentProductOutput `json:"products"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}
