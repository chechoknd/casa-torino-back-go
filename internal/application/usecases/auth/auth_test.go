package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/tests/mocks"
)

type fakeHasher struct {
	hashFn    func(string) (string, error)
	compareFn func(string, string) error
}

func (f fakeHasher) Hash(password string) (string, error) {
	return f.hashFn(password)
}

func (f fakeHasher) Compare(hash, password string) error {
	return f.compareFn(hash, password)
}

type fakeTokenIssuer struct {
	generateFn func(context.Context, uuid.UUID, string, string) (string, time.Time, error)
}

func (f fakeTokenIssuer) Generate(ctx context.Context, userID uuid.UUID, email, username string) (string, time.Time, error) {
	return f.generateFn(ctx, userID, email, username)
}

func TestRegisterCreatesUserWithHashedPassword(t *testing.T) {
	var createdUser *entities.User
	repo := &mocks.UserRepository{
		FindByEmailFn: func(context.Context, string) (*entities.User, error) {
			return nil, domainerrors.ErrNotFound
		},
		FindByUsernameFn: func(context.Context, string) (*entities.User, error) {
			return nil, domainerrors.ErrNotFound
		},
		CreateFn: func(_ context.Context, user *entities.User) error {
			createdUser = user
			return nil
		},
	}
	uc := NewUseCase(repo, fakeHasher{
		hashFn: func(password string) (string, error) {
			if password != "Password123" {
				t.Fatalf("password = %q", password)
			}
			return "hashed-password", nil
		},
		compareFn: func(string, string) error { return nil },
	}, fakeTokenIssuer{})

	out, err := uc.Register(context.Background(), dto.RegisterUserInput{
		Email:    " USER@Example.COM ",
		Username: "DemoUser",
		FullName: " Demo User ",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("Register error = %v", err)
	}
	if createdUser == nil {
		t.Fatal("expected user to be created")
	}
	if createdUser.Email != "user@example.com" || createdUser.Username != "demouser" || createdUser.FullName != "Demo User" {
		t.Fatalf("unexpected created user: %+v", createdUser)
	}
	if createdUser.PasswordHash != "hashed-password" {
		t.Fatalf("password hash = %q", createdUser.PasswordHash)
	}
	if out.Email != createdUser.Email || out.Username != createdUser.Username {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestRegisterRejectsDuplicateUsername(t *testing.T) {
	repo := &mocks.UserRepository{
		FindByEmailFn: func(context.Context, string) (*entities.User, error) {
			return nil, domainerrors.ErrNotFound
		},
		FindByUsernameFn: func(context.Context, string) (*entities.User, error) {
			return &entities.User{ID: uuid.New()}, nil
		},
		CreateFn: func(context.Context, *entities.User) error {
			t.Fatal("Create should not be called")
			return nil
		},
	}
	uc := NewUseCase(repo, fakeHasher{}, fakeTokenIssuer{})

	_, err := uc.Register(context.Background(), dto.RegisterUserInput{
		Email:    "user@example.com",
		Username: "demo",
		FullName: "Demo User",
		Password: "Password123",
	})
	if !errors.Is(err, domainerrors.ErrDuplicateUsername) {
		t.Fatalf("error = %v, want ErrDuplicateUsername", err)
	}
}

func TestLoginReturnsAccessToken(t *testing.T) {
	userID := uuid.New()
	expiresAt := time.Date(2026, 5, 8, 12, 15, 0, 0, time.UTC)
	repo := &mocks.UserRepository{
		FindByEmailFn: func(_ context.Context, email string) (*entities.User, error) {
			if email != "user@example.com" {
				t.Fatalf("email = %q", email)
			}
			return &entities.User{
				ID:           userID,
				Email:        email,
				Username:     "demo",
				FullName:     "Demo User",
				PasswordHash: "hashed-password",
				CreatedAt:    expiresAt.Add(-time.Hour),
				IsActive:     true,
			}, nil
		},
	}
	uc := NewUseCase(repo, fakeHasher{
		hashFn: func(string) (string, error) { return "", nil },
		compareFn: func(hash, password string) error {
			if hash != "hashed-password" || password != "Password123" {
				t.Fatalf("unexpected compare input")
			}
			return nil
		},
	}, fakeTokenIssuer{
		generateFn: func(_ context.Context, id uuid.UUID, email, username string) (string, time.Time, error) {
			if id != userID || email != "user@example.com" || username != "demo" {
				t.Fatalf("unexpected token input")
			}
			return "token", expiresAt, nil
		},
	})

	out, err := uc.Login(context.Background(), dto.LoginInput{
		EmailOrUsername: "USER@example.com",
		Password:        "Password123",
	})
	if err != nil {
		t.Fatalf("Login error = %v", err)
	}
	if out.AccessToken != "token" || out.TokenType != "Bearer" || !out.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("unexpected token output: %+v", out)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	repo := &mocks.UserRepository{
		FindByUsernameFn: func(context.Context, string) (*entities.User, error) {
			return &entities.User{
				ID:           uuid.New(),
				Email:        "user@example.com",
				Username:     "demo",
				PasswordHash: "hashed-password",
				IsActive:     true,
			}, nil
		},
	}
	uc := NewUseCase(repo, fakeHasher{
		hashFn: func(string) (string, error) { return "", nil },
		compareFn: func(string, string) error {
			return errors.New("password mismatch")
		},
	}, fakeTokenIssuer{})

	_, err := uc.Login(context.Background(), dto.LoginInput{
		EmailOrUsername: "demo",
		Password:        "wrong-password",
	})
	if !errors.Is(err, domainerrors.ErrInvalidCredentials) {
		t.Fatalf("error = %v, want ErrInvalidCredentials", err)
	}
}
