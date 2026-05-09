package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/interfaces/http/handlers"
	"github.com/casatorino/backend/internal/interfaces/http/routes"
)

type fakeAuthUseCase struct {
	registerFn func(context.Context, dto.RegisterUserInput) (dto.AuthUserOutput, error)
	loginFn    func(context.Context, dto.LoginInput) (dto.AuthTokenOutput, error)
	refreshFn  func(context.Context, dto.RefreshTokenInput) (dto.AuthTokenOutput, error)
	logoutFn   func(context.Context, dto.LogoutInput) error
}

func (f fakeAuthUseCase) Register(ctx context.Context, input dto.RegisterUserInput) (dto.AuthUserOutput, error) {
	return f.registerFn(ctx, input)
}

func (f fakeAuthUseCase) Login(ctx context.Context, input dto.LoginInput) (dto.AuthTokenOutput, error) {
	return f.loginFn(ctx, input)
}

func (f fakeAuthUseCase) Refresh(ctx context.Context, input dto.RefreshTokenInput) (dto.AuthTokenOutput, error) {
	return f.refreshFn(ctx, input)
}

func (f fakeAuthUseCase) Logout(ctx context.Context, input dto.LogoutInput) error {
	return f.logoutFn(ctx, input)
}

func TestAuthEndpointsSuccess(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	user := dto.AuthUserOutput{
		ID:        userID,
		Email:     "user@example.com",
		Username:  "demo",
		FullName:  "Demo User",
		CreatedAt: now,
	}
	useCase := fakeAuthUseCase{
		registerFn: func(_ context.Context, input dto.RegisterUserInput) (dto.AuthUserOutput, error) {
			if input.Email != "user@example.com" || input.Username != "demo" || input.Password != "Password123" {
				t.Fatalf("unexpected register input: %+v", input)
			}
			return user, nil
		},
		loginFn: func(_ context.Context, input dto.LoginInput) (dto.AuthTokenOutput, error) {
			if input.EmailOrUsername != "user@example.com" || input.Password != "Password123" {
				t.Fatalf("unexpected login input: %+v", input)
			}
			return dto.AuthTokenOutput{
				AccessToken: "token",
				TokenType:   "Bearer",
				ExpiresAt:   now.Add(15 * time.Minute),
				User:        user,
			}, nil
		},
		refreshFn: func(_ context.Context, input dto.RefreshTokenInput) (dto.AuthTokenOutput, error) {
			if input.RefreshToken != "refresh-token" {
				t.Fatalf("unexpected refresh input: %+v", input)
			}
			return dto.AuthTokenOutput{
				AccessToken:  "new-token",
				RefreshToken: "new-refresh-token",
				TokenType:    "Bearer",
				ExpiresAt:    now.Add(15 * time.Minute),
				User:         user,
			}, nil
		},
		logoutFn: func(_ context.Context, input dto.LogoutInput) error {
			if input.RefreshToken != "refresh-token" {
				t.Fatalf("unexpected logout input: %+v", input)
			}
			return nil
		},
	}
	router := routes.NewRouter(routes.Dependencies{
		Auth: handlers.NewAuthHandler(useCase),
	})

	tests := []struct {
		name       string
		path       string
		body       string
		wantStatus int
	}{
		{
			name:       "register",
			path:       "/auth/register",
			body:       `{"email":"user@example.com","username":"demo","full_name":"Demo User","password":"Password123"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "login",
			path:       "/auth/login",
			body:       `{"email_or_username":"user@example.com","password":"Password123"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "refresh",
			path:       "/auth/refresh",
			body:       `{"refresh_token":"refresh-token"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "logout",
			path:       "/auth/logout",
			body:       `{"refresh_token":"refresh-token"}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tc.path, bytes.NewBufferString(tc.body))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body: %s", recorder.Code, tc.wantStatus, recorder.Body.String())
			}
			assertSuccessEnvelope(t, recorder.Body.Bytes())
		})
	}
}

func TestAuthLoginInvalidCredentials(t *testing.T) {
	router := routes.NewRouter(routes.Dependencies{
		Auth: handlers.NewAuthHandler(fakeAuthUseCase{
			registerFn: func(context.Context, dto.RegisterUserInput) (dto.AuthUserOutput, error) {
				return dto.AuthUserOutput{}, nil
			},
			loginFn: func(context.Context, dto.LoginInput) (dto.AuthTokenOutput, error) {
				return dto.AuthTokenOutput{}, domainerrors.ErrInvalidCredentials
			},
			refreshFn: func(context.Context, dto.RefreshTokenInput) (dto.AuthTokenOutput, error) {
				return dto.AuthTokenOutput{}, nil
			},
			logoutFn: func(context.Context, dto.LogoutInput) error {
				return nil
			},
		}),
	})

	request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email_or_username":"user@example.com","password":"wrong"}`))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}

	var envelope struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if envelope.Code != "INVALID_CREDENTIALS" {
		t.Fatalf("code = %q", envelope.Code)
	}
}
