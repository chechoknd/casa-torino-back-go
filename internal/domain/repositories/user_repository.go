package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
}
