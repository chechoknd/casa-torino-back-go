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

type ProductHandlerUseCase interface {
	CreateProduct(ctx context.Context, input dto.CreateProductInput) (dto.ProductOutput, error)
	GetProduct(ctx context.Context, id uuid.UUID) (dto.ProductOutput, error)
	ListProducts(ctx context.Context, input dto.ListProductsInput) ([]dto.ProductOutput, error)
	UpdateProduct(ctx context.Context, input dto.UpdateProductInput) (dto.ProductOutput, error)
	DeactivateProduct(ctx context.Context, id uuid.UUID) error
}

type ProductHandler struct {
	useCase ProductHandlerUseCase
}

func NewProductHandler(useCase ProductHandlerUseCase) *ProductHandler {
	return &ProductHandler{useCase: useCase}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	var req requests.CreateProductRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	basePrice, err := requests.ParseDecimal(req.BasePrice)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	costPrice, err := requests.ParseDecimal(req.CostPrice)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.CreateProduct(r.Context(), dto.CreateProductInput{
		Name:        req.Name,
		Description: req.Description,
		ProductType: req.ProductType,
		BasePrice:   basePrice,
		CostPrice:   costPrice,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusCreated, out)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	out, err := h.useCase.ListProducts(r.Context(), dto.ListProductsInput{
		ProductType: r.URL.Query().Get("product_type"),
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.GetProduct(r.Context(), id)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	var req requests.UpdateProductRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	basePrice, err := requests.ParseDecimal(req.BasePrice)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	costPrice, err := requests.ParseDecimal(req.CostPrice)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.UpdateProduct(r.Context(), dto.UpdateProductInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		ProductType: req.ProductType,
		BasePrice:   basePrice,
		CostPrice:   costPrice,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	if err := h.useCase.DeactivateProduct(r.Context(), id); err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, map[string]string{"id": id.String()})
}
