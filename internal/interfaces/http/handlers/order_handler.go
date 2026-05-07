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

type OrderHandlerUseCase interface {
	CreateOrder(ctx context.Context, input dto.CreateOrderInput) (dto.OrderOutput, error)
	ListOrders(ctx context.Context, input dto.ListOrdersInput) ([]dto.OrderOutput, error)
	GetOrder(ctx context.Context, id uuid.UUID) (dto.OrderOutput, error)
	AddOrderItem(ctx context.Context, input dto.AddOrderItemInput) (dto.OrderOutput, error)
	UpdateOrderStatus(ctx context.Context, input dto.UpdateOrderStatusInput) (dto.OrderOutput, error)
}

type OrderHandler struct {
	useCase OrderHandlerUseCase
}

func NewOrderHandler(useCase OrderHandlerUseCase) *OrderHandler {
	return &OrderHandler{useCase: useCase}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	defer r.Body.Close()
	var req requests.CreateOrderRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	customerID, err := requests.ParseUUID(req.CustomerID)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	discount, err := requests.ParseDecimal(req.Discount)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.CreateOrder(r.Context(), dto.CreateOrderInput{
		CustomerID: customerID,
		Discount:   discount,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusCreated, out)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	var customerID *uuid.UUID
	if raw := r.URL.Query().Get("customer_id"); raw != "" {
		parsed, err := requests.ParseUUID(raw)
		if err != nil {
			responses.WriteError(w, err)
			return
		}
		customerID = &parsed
	}
	out, err := h.useCase.ListOrders(r.Context(), dto.ListOrdersInput{CustomerID: customerID})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.GetOrder(r.Context(), id)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *OrderHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	defer r.Body.Close()
	orderID, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	var req requests.AddOrderItemRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	productID, err := requests.ParseUUID(req.ProductID)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.AddOrderItem(r.Context(), dto.AddOrderItemInput{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	defer r.Body.Close()
	orderID, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	var req requests.UpdateOrderStatusRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.UpdateOrderStatus(r.Context(), dto.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  req.Status,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}
