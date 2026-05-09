package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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
const refreshTokenTTL = 7 * 24 * time.Hour

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]{3,50}$`)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenIssuer interface {
	Generate(ctx context.Context, userID uuid.UUID, email, username string) (string, time.Time, error)
}

type UseCase struct {
	users         repositories.UserRepository
	refreshTokens repositories.RefreshTokenRepository
	hasher        PasswordHasher
	tokens        TokenIssuer
}

func NewUseCase(users repositories.UserRepository, hasher PasswordHasher, tokens TokenIssuer, refreshTokens ...repositories.RefreshTokenRepository) *UseCase {
	var refreshTokenRepo repositories.RefreshTokenRepository
	if len(refreshTokens) > 0 {
		refreshTokenRepo = refreshTokens[0]
	}

	return &UseCase{
		users:         users,
		refreshTokens: refreshTokenRepo,
		hasher:        hasher,
		tokens:        tokens,
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

	refreshToken, err := uc.issueRefreshToken(ctx, user.ID)
	if err != nil {
		return dto.AuthTokenOutput{}, err
	}

	return dto.AuthTokenOutput{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
		User:         toAuthUserOutput(*user),
	}, nil
}

func (uc *UseCase) Refresh(ctx context.Context, input dto.RefreshTokenInput) (dto.AuthTokenOutput, error) {
	rawToken := strings.TrimSpace(input.RefreshToken)
	if rawToken == "" || uc.refreshTokens == nil {
		return dto.AuthTokenOutput{}, domainerrors.ErrUnauthorized
	}

	storedToken, err := uc.refreshTokens.FindByTokenHash(ctx, hashRefreshToken(rawToken))
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return dto.AuthTokenOutput{}, domainerrors.ErrUnauthorized
		}
		return dto.AuthTokenOutput{}, err
	}

	now := time.Now().UTC()
	if !storedToken.IsActive(now) {
		return dto.AuthTokenOutput{}, domainerrors.ErrUnauthorized
	}

	user, err := uc.users.FindByID(ctx, storedToken.UserID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return dto.AuthTokenOutput{}, domainerrors.ErrUnauthorized
		}
		return dto.AuthTokenOutput{}, err
	}
	if user == nil || !user.IsActive {
		return dto.AuthTokenOutput{}, domainerrors.ErrUnauthorized
	}

	if err := uc.refreshTokens.Revoke(ctx, storedToken.ID, now); err != nil {
		return dto.AuthTokenOutput{}, err
	}

	accessToken, expiresAt, err := uc.tokens.Generate(ctx, user.ID, user.Email, user.Username)
	if err != nil {
		return dto.AuthTokenOutput{}, err
	}

	newRefreshToken, err := uc.issueRefreshToken(ctx, user.ID)
	if err != nil {
		return dto.AuthTokenOutput{}, err
	}

	return dto.AuthTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
		User:         toAuthUserOutput(*user),
	}, nil
}

func (uc *UseCase) Logout(ctx context.Context, input dto.LogoutInput) error {
	rawToken := strings.TrimSpace(input.RefreshToken)
	if rawToken == "" || uc.refreshTokens == nil {
		return domainerrors.ErrUnauthorized
	}

	storedToken, err := uc.refreshTokens.FindByTokenHash(ctx, hashRefreshToken(rawToken))
	if err != nil {
		if errors.Is(err, domainerrors.ErrNotFound) {
			return domainerrors.ErrUnauthorized
		}
		return err
	}
	if storedToken.RevokedAt != nil {
		return nil
	}

	return uc.refreshTokens.Revoke(ctx, storedToken.ID, time.Now().UTC())
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

func (uc *UseCase) issueRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	if uc.refreshTokens == nil {
		return "", nil
	}

	rawToken, err := generateRefreshToken()
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	token := &entities.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: hashRefreshToken(rawToken),
		CreatedAt: now,
		ExpiresAt: now.Add(refreshTokenTTL),
	}

	if err := uc.refreshTokens.Create(ctx, token); err != nil {
		return "", err
	}

	return rawToken, nil
}

func generateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
