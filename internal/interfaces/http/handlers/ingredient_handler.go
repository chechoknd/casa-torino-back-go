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

type IngredientHandlerUseCase interface {
	CreateIngredient(ctx context.Context, input dto.CreateIngredientInput) (dto.IngredientOutput, error)
	GetIngredient(ctx context.Context, id uuid.UUID) (dto.IngredientOutput, error)
	ListIngredients(ctx context.Context) ([]dto.IngredientOutput, error)
	UpdateIngredient(ctx context.Context, input dto.UpdateIngredientInput) (dto.IngredientOutput, error)
	DeactivateIngredient(ctx context.Context, id uuid.UUID) error
}

type IngredientHandler struct {
	useCase IngredientHandlerUseCase
}

func NewIngredientHandler(useCase IngredientHandlerUseCase) *IngredientHandler {
	return &IngredientHandler{useCase: useCase}
}

func (h *IngredientHandler) Create(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	var req requests.CreateIngredientRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	averageCost, err := requests.ParseDecimal(req.AverageCost)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	stock, err := requests.ParseDecimal(req.Stock)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	minimumStock, err := requests.ParseDecimal(req.MinimumStock)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.CreateIngredient(r.Context(), dto.CreateIngredientInput{
		Name:         req.Name,
		Unit:         req.Unit,
		AverageCost:  averageCost,
		Stock:        stock,
		MinimumStock: minimumStock,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusCreated, out)
}

func (h *IngredientHandler) List(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	out, err := h.useCase.ListIngredients(r.Context())
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *IngredientHandler) Get(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.GetIngredient(r.Context(), id)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *IngredientHandler) Update(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	var req requests.UpdateIngredientRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	averageCost, err := requests.ParseDecimal(req.AverageCost)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	stock, err := requests.ParseDecimal(req.Stock)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	minimumStock, err := requests.ParseDecimal(req.MinimumStock)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.UpdateIngredient(r.Context(), dto.UpdateIngredientInput{
		ID:           id,
		Name:         req.Name,
		Unit:         req.Unit,
		AverageCost:  averageCost,
		Stock:        stock,
		MinimumStock: minimumStock,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *IngredientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	id, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	if err := h.useCase.DeactivateIngredient(r.Context(), id); err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, map[string]string{"id": id.String()})
}
