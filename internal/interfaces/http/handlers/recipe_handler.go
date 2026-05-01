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

type RecipeHandlerUseCase interface {
	CreateRecipe(ctx context.Context, input dto.CreateRecipeInput) (dto.RecipeOutput, error)
	AddRecipeItem(ctx context.Context, input dto.AddRecipeItemInput) (dto.RecipeOutput, error)
	GetRecipeByProduct(ctx context.Context, productID uuid.UUID) (dto.RecipeOutput, error)
	ListRecipes(ctx context.Context) ([]dto.RecipeOutput, error)
	GetRecipeCost(ctx context.Context, recipeID uuid.UUID) (dto.RecipeCostOutput, error)
}

type RecipeHandler struct {
	useCase RecipeHandlerUseCase
}

func NewRecipeHandler(useCase RecipeHandlerUseCase) *RecipeHandler {
	return &RecipeHandler{useCase: useCase}
}

func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	var req requests.CreateRecipeRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	productID, err := requests.ParseUUID(req.ProductID)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.CreateRecipe(r.Context(), dto.CreateRecipeInput{
		ProductID: productID,
		Name:      req.Name,
		Portions:  req.Portions,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusCreated, out)
}

func (h *RecipeHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	recipeID, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	var req requests.AddRecipeItemRequest
	if err := requests.DecodeJSON(r, &req); err != nil {
		responses.WriteError(w, err)
		return
	}
	ingredientID, err := requests.ParseUUID(req.IngredientID)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	quantity, err := requests.ParseDecimal(req.Quantity)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.AddRecipeItem(r.Context(), dto.AddRecipeItemInput{
		RecipeID:     recipeID,
		IngredientID: ingredientID,
		Quantity:     quantity,
		Unit:         req.Unit,
	})
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *RecipeHandler) GetByProduct(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	productID, err := requests.ParseUUID(chi.URLParam(r, "product_id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.GetRecipeByProduct(r.Context(), productID)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *RecipeHandler) List(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	out, err := h.useCase.ListRecipes(r.Context())
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}

func (h *RecipeHandler) GetCost(w http.ResponseWriter, r *http.Request) {
	noCache(w)
	recipeID, err := requests.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	out, err := h.useCase.GetRecipeCost(r.Context(), recipeID)
	if err != nil {
		responses.WriteError(w, err)
		return
	}
	responses.WriteJSON(w, http.StatusOK, out)
}
