package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	ingredientuc "github.com/casatorino/backend/internal/application/usecases/ingredient"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	"github.com/casatorino/backend/tests/mocks"
)

func TestCreateIngredientSuccess(t *testing.T) {
	repo := &mocks.IngredientRepository{CreateFn: func(context.Context, *entities.Ingredient) error { return nil }}
	uc := ingredientuc.NewUseCase(repo)
	_, err := uc.CreateIngredient(context.Background(), dto.CreateIngredientInput{Name: "Leche", Unit: "L"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateIngredientInvalidUnit(t *testing.T) {
	repo := &mocks.IngredientRepository{CreateFn: func(context.Context, *entities.Ingredient) error { return nil }}
	uc := ingredientuc.NewUseCase(repo)
	_, err := uc.CreateIngredient(context.Background(), dto.CreateIngredientInput{Name: "Leche", Unit: "BOX"})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestGetIngredientInactive(t *testing.T) {
	unit, _ := valueobjects.NewUnit("KG")
	repo := &mocks.IngredientRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Ingredient, error) {
			return &entities.Ingredient{ID: uuid.New(), Unit: unit, IsActive: false}, nil
		},
	}
	uc := ingredientuc.NewUseCase(repo)
	_, err := uc.GetIngredient(context.Background(), uuid.New())
	if !errors.Is(err, domainerrors.ErrInactive) {
		t.Fatalf("expected inactive error, got %v", err)
	}
}
