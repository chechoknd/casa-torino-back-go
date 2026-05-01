package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/entities"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	Update(ctx context.Context, order *entities.Order) error
	AddItem(ctx context.Context, orderID uuid.UUID, item *entities.OrderItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	List(ctx context.Context) ([]entities.Order, error)
	ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]entities.Order, error)
}
