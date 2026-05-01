package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/interfaces/http/requests"
	"github.com/casatorino/backend/internal/interfaces/http/responses"
)

type CustomerHandlerUseCase interface {
	CreateCustomer(ctx context.Context, input dto.CreateCustomerInput) (dto.CustomerOutput, error)
	GetCustomer(ctx context.Context, id uuid.UUID) (dto.CustomerOutput, error)
	ListCustomers(ctx context.Context) ([]dto.CustomerOutput, error)
	UpdateCustomer(ctx context.Context, input dto.UpdateCustomerInput) (dto.CustomerOutput, error)
	DeactivateCustomer(ctx context.Context, id uuid.UUID) error
}

type CustomerHandler struct {
	useCase CustomerHandlerUseCase
}

func NewCustomerHandler(useCase CustomerHandlerUseCase) *CustomerHandler {
	return &CustomerHandler{useCase: useCase}
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	var req requests.CreateCustomerRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}

	out, err := h.useCase.CreateCustomer(r.Context(), dto.CreateCustomerInput{
		FullName:     req.FullName,
		Phone:        req.Phone,
		Email:        req.Email,
		CustomerType: req.CustomerType,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusCreated, out)
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	out, err := h.useCase.ListCustomers(r.Context())
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.GetCustomer(r.Context(), id)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	var req requests.UpdateCustomerRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.UpdateCustomer(r.Context(), dto.UpdateCustomerInput{
		ID:           id,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Email:        req.Email,
		CustomerType: req.CustomerType,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	if err := h.useCase.DeactivateCustomer(r.Context(), id); err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, map[string]string{"id": id.String()})
}
