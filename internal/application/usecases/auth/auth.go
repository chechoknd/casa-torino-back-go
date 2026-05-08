package auth

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/repositories"
)

const minPasswordLength = 8

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{3,50}$`)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenIssuer interface {
	Generate(ctx context.Context, userID uuid.UUID, email, username string) (string, time.Time, error)
}

type UseCase struct {
	users  repositories.UserRepository
	hasher PasswordHasher
	tokens TokenIssuer
}

func NewUseCase(users repositories.UserRepository, hasher PasswordHasher, tokens TokenIssuer) *UseCase {
	return &UseCase{
		users:  users,
		hasher: hasher,
		tokens: tokens,
	}
}

func (uc *UseCase) Register(ctx context.Context, input dto.RegisterUserInput) (dto.AuthUserOutput, error) {
	email := normalizeEmail(input.Email)
	username := normalizeUsername(input.Username)
	fullName := strings.TrimSpace(input.FullName)

	if err := validateRegistration(email, username, fullName, input.Password); err != nil {
		return dto.AuthUserOutput{}, err
	}

	existing, err := uc.users.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return dto.AuthUserOutput{}, err
	}
	if existing != nil {
		return dto.AuthUserOutput{}, domainerrors.ErrDuplicateEmail
	}

	existing, err = uc.users.FindByUsername(ctx, username)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return dto.AuthUserOutput{}, err
	}
	if existing != nil {
		return dto.AuthUserOutput{}, domainerrors.ErrDuplicateUsername
	}

	passwordHash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return dto.AuthUserOutput{}, err
	}

	now := time.Now().UTC()
	user := &entities.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		FullName:     fullName,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}

	if err := uc.users.Create(ctx, user); err != nil {
		return dto.AuthUserOutput{}, err
	}

	return toAuthUserOutput(*user), nil
}

func (uc *UseCase) Login(ctx context.Context, input dto.LoginInput) (dto.AuthTokenOutput, error) {
	identifier := strings.TrimSpace(input.EmailOrUsername)
	if identifier == "" || strings.TrimSpace(input.Password) == "" {
		return dto.AuthTokenOutput{}, domainerrors.ErrInvalidCredentials
	}

	user, err := uc.findUserByIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return dto.AuthTokenOutput{}, domainerrors.ErrInvalidCredentials
		}
		return dto.AuthTokenOutput{}, err
	}
	if user == nil || !user.IsActive {
		return dto.AuthTokenOutput{}, domainerrors.ErrInvalidCredentials
	}

	if err := uc.hasher.Compare(user.PasswordHash, input.Password); err != nil {
		return dto.AuthTokenOutput{}, domainerrors.ErrInvalidCredentials
	}

	token, expiresAt, err := uc.tokens.Generate(ctx, user.ID, user.Email, user.Username)
	if err != nil {
		return dto.AuthTokenOutput{}, err
	}

	return dto.AuthTokenOutput{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		User:        toAuthUserOutput(*user),
	}, nil
}

func (uc *UseCase) findUserByIdentifier(ctx context.Context, identifier string) (*entities.User, error) {
	if strings.Contains(identifier, "@") {
		return uc.users.FindByEmail(ctx, normalizeEmail(identifier))
	}

	return uc.users.FindByUsername(ctx, normalizeUsername(identifier))
}

func validateRegistration(email, username, fullName, password string) error {
	if email == "" || username == "" || fullName == "" || strings.TrimSpace(password) == "" {
		return domainerrors.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return domainerrors.ErrInvalidInput
	}
	if !usernamePattern.MatchString(username) {
		return domainerrors.ErrInvalidInput
	}
	if len(password) < minPasswordLength {
		return domainerrors.ErrInvalidInput
	}

	return nil
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func normalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func toAuthUserOutput(user entities.User) dto.AuthUserOutput {
	return dto.AuthUserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
	}
}
