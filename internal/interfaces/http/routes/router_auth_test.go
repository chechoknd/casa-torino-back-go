package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	"github.com/casatorino/backend/internal/interfaces/http/handlers"
	"github.com/casatorino/backend/internal/interfaces/http/middleware"
)

type rejectingVerifier struct{}

func (rejectingVerifier) Verify(context.Context, string) (middleware.TokenClaims, error) {
	return middleware.TokenClaims{}, nil
}

type customerVerifier struct{}

func (customerVerifier) Verify(context.Context, string) (middleware.TokenClaims, error) {
	return middleware.TokenClaims{
		Role: valueobjects.UserRoleCustomer,
	}, nil
}

type fakeAuthUseCase struct{}

func (fakeAuthUseCase) Register(context.Context, dto.RegisterUserInput) (dto.AuthUserOutput, error) {
	return dto.AuthUserOutput{}, nil
}

func (fakeAuthUseCase) Login(context.Context, dto.LoginInput) (dto.AuthTokenOutput, error) {
	return dto.AuthTokenOutput{}, nil
}

func (fakeAuthUseCase) Refresh(context.Context, dto.RefreshTokenInput) (dto.AuthTokenOutput, error) {
	return dto.AuthTokenOutput{}, nil
}

func (fakeAuthUseCase) Logout(context.Context, dto.LogoutInput) error {
	return nil
}

func TestRouterProtectsBusinessRoutesWhenVerifierConfigured(t *testing.T) {
	router := NewRouter(Dependencies{
		TokenVerifier: rejectingVerifier{},
	})

	request := httptest.NewRequest(http.MethodGet, "/products/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestRouterRequiresAdminRoleForBusinessRoutes(t *testing.T) {
	router := NewRouter(Dependencies{
		TokenVerifier: customerVerifier{},
	})

	request := httptest.NewRequest(http.MethodGet, "/products/", nil)
	request.Header.Set("Authorization", "Bearer customer-token")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func TestRouterKeepsAuthRoutesPublic(t *testing.T) {
	router := NewRouter(Dependencies{
		Auth:          handlers.NewAuthHandler(fakeAuthUseCase{}),
		TokenVerifier: rejectingVerifier{},
	})

	request := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d from public auth handler validation", recorder.Code, http.StatusBadRequest)
	}
}

func TestRouterRejectsTooLargeRequestBody(t *testing.T) {
	router := NewRouter(Dependencies{
		Auth: handlers.NewAuthHandler(fakeAuthUseCase{}),
	})

	body := `{"email_or_username":"` + strings.Repeat("a", 1<<20) + `","password":"password123"}`
	request := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestRouterHealthIsPublic(t *testing.T) {
	router := NewRouter(Dependencies{
		TokenVerifier: rejectingVerifier{},
	})

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if body := recorder.Body.String(); body != `{"status":"ok"}` {
		t.Fatalf("body = %q, want health status", body)
	}
}
