package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/interfaces/http/responses"
)

type AuthenticatedUser struct {
	ID        uuid.UUID
	Email     string
	Username  string
	ExpiresAt time.Time
}

type TokenClaims struct {
	UserID    uuid.UUID
	Email     string
	Username  string
	ExpiresAt time.Time
}

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (TokenClaims, error)
}

type authContextKey struct{}

func JWTAuth(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := bearerToken(r.Header.Get("Authorization"))
			if err != nil {
				responses.WriteError(w, err)
				return
			}

			claims, err := verifier.Verify(r.Context(), token)
			if err != nil {
				responses.WriteError(w, domainerrors.ErrUnauthorized)
				return
			}

			user := AuthenticatedUser{
				ID:        claims.UserID,
				Email:     claims.Email,
				Username:  claims.Username,
				ExpiresAt: claims.ExpiresAt,
			}
			ctx := context.WithValue(r.Context(), authContextKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) (AuthenticatedUser, bool) {
	user, ok := ctx.Value(authContextKey{}).(AuthenticatedUser)
	return user, ok
}

func bearerToken(header string) (string, error) {
	scheme, token, ok := strings.Cut(strings.TrimSpace(header), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
		return "", domainerrors.ErrUnauthorized
	}
	if strings.Contains(token, " ") {
		return "", domainerrors.ErrUnauthorized
	}

	return token, nil
}
