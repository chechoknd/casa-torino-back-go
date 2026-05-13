package security

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type TokenClaims struct {
	UserID    uuid.UUID
	Email     string
	Username  string
	Role      valueobjects.UserRole
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type JWTManager struct {
	secret    []byte
	expiresIn time.Duration
	now       func() time.Time
}

func NewJWTManager(secret string, expiresIn time.Duration) *JWTManager {
	return &JWTManager{
		secret:    []byte(secret),
		expiresIn: expiresIn,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

func NewJWTManagerWithClock(secret string, expiresIn time.Duration, now func() time.Time) *JWTManager {
	manager := NewJWTManager(secret, expiresIn)
	manager.now = now
	return manager
}

func (m *JWTManager) Generate(_ context.Context, userID uuid.UUID, email, username string, role valueobjects.UserRole) (string, time.Time, error) {
	now := m.now()
	expiresAt := now.Add(m.expiresIn)

	claims := jwtClaims{
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

func (m *JWTManager) Verify(_ context.Context, token string) (TokenClaims, error) {
	claims := jwtClaims{}
	parser := jwt.NewParser(
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(m.now),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	parsedToken, err := parser.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, domainerrors.ErrUnauthorized
		}
		return m.secret, nil
	})
	if err != nil || !parsedToken.Valid {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}

	if claims.Subject == "" || claims.Email == "" || claims.Username == "" || claims.Role == "" || claims.ExpiresAt == nil {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}

	role, err := valueobjects.NewUserRole(string(claims.Role))
	if err != nil {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}

	issuedAt := time.Time{}
	if claims.IssuedAt != nil {
		issuedAt = claims.IssuedAt.Time.UTC()
	}

	return TokenClaims{
		UserID:    userID,
		Email:     claims.Email,
		Username:  claims.Username,
		Role:      role,
		IssuedAt:  issuedAt,
		ExpiresAt: claims.ExpiresAt.Time.UTC(),
	}, nil
}

type jwtClaims struct {
	Email    string                `json:"email"`
	Username string                `json:"username"`
	Role     valueobjects.UserRole `json:"role"`
	jwt.RegisteredClaims
}
