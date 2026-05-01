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

type PaymentOutput struct {
	ID        uuid.UUID       `json:"id"`
	OrderID   uuid.UUID       `json:"order_id"`
	Amount    decimal.Decimal `json:"amount"`
	Method    string          `json:"method"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
