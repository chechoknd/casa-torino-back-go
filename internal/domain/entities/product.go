package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	ProductType valueobjects.ProductType
	BasePrice   decimal.Decimal
	CostPrice   decimal.Decimal
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
