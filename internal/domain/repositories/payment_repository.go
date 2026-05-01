package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/entities"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entities.Payment) error
	Update(ctx context.Context, payment *entities.Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)
	ListByOrderID(ctx context.Context, orderID uuid.UUID) ([]entities.Payment, error)
}
