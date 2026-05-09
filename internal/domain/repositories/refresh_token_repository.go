package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/entities"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error
}
