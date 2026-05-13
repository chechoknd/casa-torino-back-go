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
	"github.com/casatorino/backend/internal/domain/valueobjects"
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
	generateFn func(context.Context, uuid.UUID, string, string, valueobjects.UserRole) (string, time.Time, error)
}

func (f fakeTokenIssuer) Generate(ctx context.Context, userID uuid.UUID, email, username string, role valueobjects.UserRole) (string, time.Time, error) {
	return f.generateFn(ctx, userID, email, username, role)
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
	if createdUser.Role != valueobjects.UserRoleCustomer {
		t.Fatalf("role = %q, want CUSTOMER", createdUser.Role)
	}
	if out.Email != createdUser.Email || out.Username != createdUser.Username || out.Role != valueobjects.UserRoleCustomer {
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
				Role:         valueobjects.UserRoleAdmin,
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
		generateFn: func(_ context.Context, id uuid.UUID, email, username string, role valueobjects.UserRole) (string, time.Time, error) {
			if id != userID || email != "user@example.com" || username != "demo" || role != valueobjects.UserRoleAdmin {
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
	if out.User.Role != valueobjects.UserRoleAdmin {
		t.Fatalf("role = %q, want ADMIN", out.User.Role)
	}
}

func TestLoginCreatesRefreshToken(t *testing.T) {
	userID := uuid.New()
	expiresAt := time.Date(2026, 5, 8, 12, 15, 0, 0, time.UTC)
	var createdToken *entities.RefreshToken

	repo := &mocks.UserRepository{
		FindByEmailFn: func(context.Context, string) (*entities.User, error) {
			return &entities.User{
				ID:           userID,
				Email:        "user@example.com",
				Username:     "demo",
				FullName:     "Demo User",
				Role:         valueobjects.UserRoleCustomer,
				PasswordHash: "hashed-password",
				IsActive:     true,
			}, nil
		},
	}
	refreshRepo := &mocks.RefreshTokenRepository{
		CreateFn: func(_ context.Context, token *entities.RefreshToken) error {
			createdToken = token
			return nil
		},
	}
	uc := NewUseCase(repo, fakeHasher{
		hashFn:    func(string) (string, error) { return "", nil },
		compareFn: func(string, string) error { return nil },
	}, fakeTokenIssuer{
		generateFn: func(context.Context, uuid.UUID, string, string, valueobjects.UserRole) (string, time.Time, error) {
			return "access-token", expiresAt, nil
		},
	}, refreshRepo)

	out, err := uc.Login(context.Background(), dto.LoginInput{
		EmailOrUsername: "user@example.com",
		Password:        "Password123",
	})
	if err != nil {
		t.Fatalf("Login error = %v", err)
	}
	if out.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
	if createdToken == nil {
		t.Fatal("expected refresh token to be stored")
	}
	if createdToken.UserID != userID || createdToken.TokenHash == "" || createdToken.TokenHash == out.RefreshToken {
		t.Fatalf("unexpected stored refresh token: %+v", createdToken)
	}
}

func TestRefreshRotatesRefreshToken(t *testing.T) {
	userID := uuid.New()
	tokenID := uuid.New()
	oldRefreshToken := "old-refresh-token"
	expiresAt := time.Date(2026, 5, 8, 12, 15, 0, 0, time.UTC)
	var revokedTokenID uuid.UUID
	var createdToken *entities.RefreshToken

	userRepo := &mocks.UserRepository{
		FindByIDFn: func(_ context.Context, id uuid.UUID) (*entities.User, error) {
			if id != userID {
				t.Fatalf("user id = %s, want %s", id, userID)
			}
			return &entities.User{
				ID:       userID,
				Email:    "user@example.com",
				Username: "demo",
				FullName: "Demo User",
				Role:     valueobjects.UserRoleCustomer,
				IsActive: true,
			}, nil
		},
	}
	refreshRepo := &mocks.RefreshTokenRepository{
		FindByTokenHashFn: func(_ context.Context, tokenHash string) (*entities.RefreshToken, error) {
			if tokenHash != hashRefreshToken(oldRefreshToken) {
				t.Fatalf("token hash = %q", tokenHash)
			}
			return &entities.RefreshToken{
				ID:        tokenID,
				UserID:    userID,
				TokenHash: tokenHash,
				CreatedAt: time.Now().UTC().Add(-time.Hour),
				ExpiresAt: time.Now().UTC().Add(time.Hour),
			}, nil
		},
		RevokeFn: func(_ context.Context, id uuid.UUID, revokedAt time.Time) error {
			revokedTokenID = id
			if revokedAt.IsZero() {
				t.Fatal("expected revoked timestamp")
			}
			return nil
		},
		CreateFn: func(_ context.Context, token *entities.RefreshToken) error {
			createdToken = token
			return nil
		},
	}
	uc := NewUseCase(userRepo, fakeHasher{}, fakeTokenIssuer{
		generateFn: func(_ context.Context, id uuid.UUID, email, username string, role valueobjects.UserRole) (string, time.Time, error) {
			if id != userID || email != "user@example.com" || username != "demo" || role != valueobjects.UserRoleCustomer {
				t.Fatalf("unexpected access token input")
			}
			return "new-access-token", expiresAt, nil
		},
	}, refreshRepo)

	out, err := uc.Refresh(context.Background(), dto.RefreshTokenInput{RefreshToken: oldRefreshToken})
	if err != nil {
		t.Fatalf("Refresh error = %v", err)
	}
	if revokedTokenID != tokenID {
		t.Fatalf("revoked token id = %s, want %s", revokedTokenID, tokenID)
	}
	if createdToken == nil {
		t.Fatal("expected rotated refresh token to be stored")
	}
	if out.AccessToken != "new-access-token" || out.RefreshToken == "" || out.RefreshToken == oldRefreshToken {
		t.Fatalf("unexpected refresh output: %+v", out)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	repo := &mocks.UserRepository{
		FindByUsernameFn: func(context.Context, string) (*entities.User, error) {
			return &entities.User{
				ID:           uuid.New(),
				Email:        "user@example.com",
				Username:     "demo",
				Role:         valueobjects.UserRoleCustomer,
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
