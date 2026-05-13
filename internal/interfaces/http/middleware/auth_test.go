package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type fakeVerifier struct {
	verifyFn func(context.Context, string) (TokenClaims, error)
}

func (f fakeVerifier) Verify(ctx context.Context, token string) (TokenClaims, error) {
	return f.verifyFn(ctx, token)
}

func TestJWTAuthInjectsAuthenticatedUser(t *testing.T) {
	userID := uuid.New()
	expiresAt := time.Date(2026, 5, 8, 12, 15, 0, 0, time.UTC)
	handler := JWTAuth(fakeVerifier{
		verifyFn: func(_ context.Context, token string) (TokenClaims, error) {
			if token != "valid-token" {
				t.Fatalf("token = %q", token)
			}
			return TokenClaims{
				UserID:    userID,
				Email:     "user@example.com",
				Username:  "demo",
				Role:      valueobjects.UserRoleAdmin,
				ExpiresAt: expiresAt,
			}, nil
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok {
			t.Fatal("expected authenticated user in context")
		}
		if user.ID != userID || user.Email != "user@example.com" || user.Username != "demo" || user.Role != valueobjects.UserRoleAdmin {
			t.Fatalf("unexpected authenticated user: %+v", user)
		}
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/customers", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestJWTAuthRejectsMissingToken(t *testing.T) {
	handler := JWTAuth(fakeVerifier{
		verifyFn: func(context.Context, string) (TokenClaims, error) {
			t.Fatal("Verify should not be called")
			return TokenClaims{}, nil
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	request := httptest.NewRequest(http.MethodGet, "/customers", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestJWTAuthRejectsInvalidToken(t *testing.T) {
	handler := JWTAuth(fakeVerifier{
		verifyFn: func(context.Context, string) (TokenClaims, error) {
			return TokenClaims{}, errors.New("invalid token")
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	request := httptest.NewRequest(http.MethodGet, "/customers", nil)
	request.Header.Set("Authorization", "Bearer invalid-token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestBearerTokenRejectsMalformedHeader(t *testing.T) {
	_, err := bearerToken("Basic token")
	if !errors.Is(err, domainerrors.ErrUnauthorized) {
		t.Fatalf("error = %v, want ErrUnauthorized", err)
	}
}

func TestRequireRoleAllowsMatchingRole(t *testing.T) {
	userID := uuid.New()
	handler := JWTAuth(fakeVerifier{
		verifyFn: func(context.Context, string) (TokenClaims, error) {
			return TokenClaims{
				UserID:   userID,
				Email:    "admin@example.com",
				Username: "admin",
				Role:     valueobjects.UserRoleAdmin,
			}, nil
		},
	})(RequireRole(valueobjects.UserRoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})))

	request := httptest.NewRequest(http.MethodGet, "/products", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}

func TestRequireRoleRejectsDifferentRole(t *testing.T) {
	handler := JWTAuth(fakeVerifier{
		verifyFn: func(context.Context, string) (TokenClaims, error) {
			return TokenClaims{
				UserID:   uuid.New(),
				Email:    "customer@example.com",
				Username: "customer",
				Role:     valueobjects.UserRoleCustomer,
			}, nil
		},
	})(RequireRole(valueobjects.UserRoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})))

	request := httptest.NewRequest(http.MethodGet, "/products", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}
