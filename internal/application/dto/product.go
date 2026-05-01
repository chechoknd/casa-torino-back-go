package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductInput struct {
	Name        string
	Description string
	ProductType string
	BasePrice   decimal.Decimal
	CostPrice   decimal.Decimal
}

type UpdateProductInput struct {
	ID          uuid.UUID
	Name        string
	Description string
	ProductType string
	BasePrice   decimal.Decimal
	CostPrice   decimal.Decimal
}

type ListProductsInput struct {
	ProductType string
}

type ProductOutput struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ProductType string          `json:"product_type"`
	BasePrice   decimal.Decimal `json:"base_price"`
	CostPrice   decimal.Decimal `json:"cost_price"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
