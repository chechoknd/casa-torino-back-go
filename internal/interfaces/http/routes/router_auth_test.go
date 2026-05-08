package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/interfaces/http/handlers"
	"github.com/casatorino/backend/internal/interfaces/http/middleware"
)

type rejectingVerifier struct{}

func (rejectingVerifier) Verify(context.Context, string) (middleware.TokenClaims, error) {
	return middleware.TokenClaims{}, nil
}

type fakeAuthUseCase struct{}

func (fakeAuthUseCase) Register(context.Context, dto.RegisterUserInput) (dto.AuthUserOutput, error) {
	return dto.AuthUserOutput{}, nil
}

func (fakeAuthUseCase) Login(context.Context, dto.LoginInput) (dto.AuthTokenOutput, error) {
	return dto.AuthTokenOutput{}, nil
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
