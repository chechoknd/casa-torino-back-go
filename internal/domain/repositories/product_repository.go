package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/entities"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	Update(ctx context.Context, product *entities.Product) error
	Deactivate(ctx context.Context, id uuid.UUID, updatedAt time.Time) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	ListActive(ctx context.Context) ([]entities.Product, error)
}
