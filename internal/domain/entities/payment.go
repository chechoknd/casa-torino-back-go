package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type Payment struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	Amount    decimal.Decimal
	Method    valueobjects.PaymentMethod
	Status    valueobjects.PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
