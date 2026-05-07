package ingredient

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/repositories"
	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type UseCase struct {
	ingredients repositories.IngredientRepository
}

func NewUseCase(ingredients repositories.IngredientRepository) *UseCase {
	return &UseCase{ingredients: ingredients}
}

func (uc *UseCase) CreateIngredient(ctx context.Context, input dto.CreateIngredientInput) (dto.IngredientOutput, error) {
	if strings.TrimSpace(input.Name) == "" {
		return dto.IngredientOutput{}, domainerrors.ErrInvalidInput
	}

	unit, err := valueobjects.NewUnit(input.Unit)
	if err != nil {
		return dto.IngredientOutput{}, err
	}

	now := time.Now().UTC()
	ingredient := &entities.Ingredient{
		ID:           uuid.New(),
		Name:         strings.TrimSpace(input.Name),
		Unit:         unit,
		AverageCost:  input.AverageCost,
		Stock:        input.Stock,
		MinimumStock: input.MinimumStock,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}

	if err := uc.ingredients.Create(ctx, ingredient); err != nil {
		return dto.IngredientOutput{}, err
	}

	return toIngredientOutput(*ingredient), nil
}

func (uc *UseCase) GetIngredient(ctx context.Context, id uuid.UUID) (dto.IngredientOutput, error) {
	ingredient, err := uc.ingredients.FindByID(ctx, id)
	if err != nil {
		return dto.IngredientOutput{}, err
	}
	if ingredient == nil || !ingredient.IsActive {
		return dto.IngredientOutput{}, domainerrors.ErrInactive
	}

	return toIngredientOutput(*ingredient), nil
}

func (uc *UseCase) ListIngredients(ctx context.Context) ([]dto.IngredientOutput, error) {
	ingredients, err := uc.ingredients.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]dto.IngredientOutput, 0, len(ingredients))
	for _, ingredient := range ingredients {
		if !ingredient.IsActive {
			continue
		}
		output = append(output, toIngredientOutput(ingredient))
	}

	return output, nil
}

func (uc *UseCase) UpdateIngredient(ctx context.Context, input dto.UpdateIngredientInput) (dto.IngredientOutput, error) {
	ingredient, err := uc.ingredients.FindByID(ctx, input.ID)
	if err != nil {
		return dto.IngredientOutput{}, err
	}
	if !ingredient.IsActive {
		return dto.IngredientOutput{}, domainerrors.ErrInactive
	}
	if strings.TrimSpace(input.Name) == "" {
		return dto.IngredientOutput{}, domainerrors.ErrInvalidInput
	}

	unit, err := valueobjects.NewUnit(input.Unit)
	if err != nil {
		return dto.IngredientOutput{}, err
	}

	ingredient.Name = strings.TrimSpace(input.Name)
	ingredient.Unit = unit
	ingredient.AverageCost = input.AverageCost
	ingredient.Stock = input.Stock
	ingredient.MinimumStock = input.MinimumStock
	ingredient.UpdatedAt = time.Now().UTC()

	if err := uc.ingredients.Update(ctx, ingredient); err != nil {
		return dto.IngredientOutput{}, err
	}

	return toIngredientOutput(*ingredient), nil
}

func (uc *UseCase) DeactivateIngredient(ctx context.Context, id uuid.UUID) error {
	ingredient, err := uc.ingredients.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if ingredient == nil || !ingredient.IsActive {
		return domainerrors.ErrInactive
	}

	return uc.ingredients.Deactivate(ctx, id, time.Now().UTC())
}

func toIngredientOutput(ingredient entities.Ingredient) dto.IngredientOutput {
	return dto.IngredientOutput{
		ID:           ingredient.ID,
		Name:         ingredient.Name,
		Unit:         string(ingredient.Unit),
		AverageCost:  ingredient.AverageCost,
		Stock:        ingredient.Stock,
		MinimumStock: ingredient.MinimumStock,
		CreatedAt:    ingredient.CreatedAt,
		UpdatedAt:    ingredient.UpdatedAt,
		IsActive:     ingredient.IsActive,
	}
}
