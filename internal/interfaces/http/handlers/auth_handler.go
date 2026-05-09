package handlers

import (
	"context"
	"net/http"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/interfaces/http/requests"
	"github.com/casatorino/backend/internal/interfaces/http/responses"
)

type AuthHandlerUseCase interface {
	Register(ctx context.Context, input dto.RegisterUserInput) (dto.AuthUserOutput, error)
	Login(ctx context.Context, input dto.LoginInput) (dto.AuthTokenOutput, error)
	Refresh(ctx context.Context, input dto.RefreshTokenInput) (dto.AuthTokenOutput, error)
	Logout(ctx context.Context, input dto.LogoutInput) error
}

type AuthHandler struct {
	useCase AuthHandlerUseCase
}

func NewAuthHandler(useCase AuthHandlerUseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	defer r.Body.Close()

	var req requests.RegisterUserRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}

	out, err := h.useCase.Register(r.Context(), dto.RegisterUserInput{
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
		Password: req.Password,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusCreated, out)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	defer r.Body.Close()

	var req requests.LoginRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	identifier := req.EmailOrUsername
	if identifier == "" {
		identifier = req.Email
	}
	if identifier == "" {
		identifier = req.Username
	}

	out, err := h.useCase.Login(r.Context(), dto.LoginInput{
		EmailOrUsername: identifier,
		Password:        req.Password,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	defer r.Body.Close()

	var req requests.RefreshTokenRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}

	out, err := h.useCase.Refresh(r.Context(), dto.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	defer r.Body.Close()

	var req requests.LogoutRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}

	if err := h.useCase.Logout(r.Context(), dto.LogoutInput{
		RefreshToken: req.RefreshToken,
	}); err != nil {
		responses.WriteError(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
