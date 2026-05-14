package handlers

import (
	"context"
	"net/http"

	"github.com/casatorino/backend/internal/application/dto"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	appmiddleware "github.com/casatorino/backend/internal/interfaces/http/middleware"
	"github.com/casatorino/backend/internal/interfaces/http/responses"
)

type CustomerPanelUseCase interface {
	GetProfile(ctx context.Context, userEmail string) (dto.CustomerOutput, error)
	ListOrders(ctx context.Context, userEmail string) ([]dto.OrderOutput, error)
}

type CustomerPanelHandler struct {
	useCase CustomerPanelUseCase
}

func NewCustomerPanelHandler(useCase CustomerPanelUseCase) *CustomerPanelHandler {
	return &CustomerPanelHandler{useCase: useCase}
}

func (h *CustomerPanelHandler) Profile(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	user, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		responses.WriteError(w, domainerrors.ErrUnauthorized)
		return
	}

	out, err := h.useCase.GetProfile(r.Context(), user.Email)
	if err != nil {
		responses.WriteError(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *CustomerPanelHandler) Orders(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	user, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		responses.WriteError(w, domainerrors.ErrUnauthorized)
		return
	}

	out, err := h.useCase.ListOrders(r.Context(), user.Email)
	if err != nil {
		responses.WriteError(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, out)
}
