package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/entities"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entities.Customer) error
	Update(ctx context.Context, customer *entities.Customer) error
	Deactivate(ctx context.Context, id uuid.UUID, updatedAt time.Time) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error)
	FindByEmail(ctx context.Context, email string) (*entities.Customer, error)
	List(ctx context.Context) ([]entities.Customer, error)
}
