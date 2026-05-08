package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casatorino/backend/internal/domain/entities"
	domainrepositories "github.com/casatorino/backend/internal/domain/repositories"
)

var _ domainrepositories.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

type RefreshTokenRepository struct {
	conn *pgxpool.Pool
}

func NewRefreshTokenRepository(conn *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{conn: conn}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	row := r.conn.QueryRow(ctx, `
		INSERT INTO refresh_tokens (
			id,
			user_id,
			token_hash,
			created_at,
			expires_at,
			revoked_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id, user_id, token_hash, created_at, expires_at, revoked_at
	`,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.CreatedAt,
		token.ExpiresAt,
		token.RevokedAt,
	)

	return scanRefreshToken(row, token)
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error) {
	var token entities.RefreshToken
	row := r.conn.QueryRow(ctx, `
		SELECT id, user_id, token_hash, created_at, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`, tokenHash)

	if err := scanRefreshToken(row, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	result, err := r.conn.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE id = $1
	`, id, revokedAt)
	if err != nil {
		return mapError(err)
	}
	if result.RowsAffected() == 0 {
		return mapError(pgx.ErrNoRows)
	}

	return nil
}

type refreshTokenScanner interface {
	Scan(dest ...any) error
}

func scanRefreshToken(row refreshTokenScanner, token *entities.RefreshToken) error {
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.CreatedAt,
		&token.ExpiresAt,
		&token.RevokedAt,
	)
	if err != nil {
		return mapError(err)
	}

	return nil
}

var _ refreshTokenScanner = pgx.Row(nil)
