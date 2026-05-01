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

type PaymentHandlerUseCase interface {
	CreatePayment(ctx context.Context, input dto.CreatePaymentInput) (dto.PaymentOutput, error)
	GetPaymentsByOrder(ctx context.Context, orderID uuid.UUID) ([]dto.PaymentOutput, error)
	ListPayments(ctx context.Context) ([]dto.PaymentOutput, error)
	UpdatePaymentStatus(ctx context.Context, input dto.UpdatePaymentStatusInput) (dto.PaymentOutput, error)
}

type PaymentHandler struct {
	useCase PaymentHandlerUseCase
}

func NewPaymentHandler(useCase PaymentHandlerUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	var req requests.CreatePaymentRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	orderID, err := requests.ParseUUID(req.OrderID)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	amount, err := requests.ParseDecimal(req.Amount)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.CreatePayment(r.Context(), dto.CreatePaymentInput{
		OrderID: orderID,
		Amount:  amount,
		Method:  req.Method,
		Status:  req.Status,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusCreated, out)
}

func (h *PaymentHandler) GetByOrder(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	orderID, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.GetPaymentsByOrder(r.Context(), orderID)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	out, err := h.useCase.ListPayments(r.Context())
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *PaymentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	paymentID, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	var req requests.UpdatePaymentStatusRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.UpdatePaymentStatus(r.Context(), dto.UpdatePaymentStatusInput{
		PaymentID: paymentID,
		Status:    req.Status,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}
