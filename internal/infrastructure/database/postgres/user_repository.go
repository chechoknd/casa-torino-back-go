package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casatorino/backend/internal/domain/entities"
	domainrepositories "github.com/casatorino/backend/internal/domain/repositories"
)

var _ domainrepositories.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{conn: conn}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	row := r.conn.QueryRow(ctx, `
		INSERT INTO users (
			id,
			email,
			username,
			full_name,
			password_hash,
			created_at,
			updated_at,
			is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING id, email, username, full_name, password_hash, created_at, updated_at, is_active
	`,
		user.ID,
		user.Email,
		user.Username,
		user.FullName,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
		user.IsActive,
	)

	return scanUser(row, user)
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	row := r.conn.QueryRow(ctx, `
		SELECT id, email, username, full_name, password_hash, created_at, updated_at, is_active
		FROM users
		WHERE id = $1
	`, id)

	if err := scanUser(row, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	row := r.conn.QueryRow(ctx, `
		SELECT id, email, username, full_name, password_hash, created_at, updated_at, is_active
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`, email)

	if err := scanUser(row, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	row := r.conn.QueryRow(ctx, `
		SELECT id, email, username, full_name, password_hash, created_at, updated_at, is_active
		FROM users
		WHERE LOWER(username) = LOWER($1)
	`, username)

	if err := scanUser(row, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(row userScanner, user *entities.User) error {
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.FullName,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		return mapError(err)
	}

	return nil
}

var _ userScanner = pgx.Row(nil)
