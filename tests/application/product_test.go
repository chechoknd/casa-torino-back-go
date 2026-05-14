package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/application/dto"
	productuc "github.com/casatorino/backend/internal/application/usecases/product"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	"github.com/casatorino/backend/tests/mocks"
)

func TestCreateProductSuccess(t *testing.T) {
	repo := &mocks.ProductRepository{CreateFn: func(context.Context, *entities.Product) error { return nil }}
	uc := productuc.NewUseCase(repo)
	out, err := uc.CreateProduct(context.Background(), dto.CreateProductInput{
		Name: "Menu", ProductType: "LUNCH", BasePrice: decimal.RequireFromString("10000"), ImageURL: "https://cdn.example.com/menu.jpg", IsPublic: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ProductType != "LUNCH" {
		t.Fatalf("unexpected type: %s", out.ProductType)
	}
	if out.ImageURL != "https://cdn.example.com/menu.jpg" || !out.IsPublic {
		t.Fatalf("unexpected catalog fields: %+v", out)
	}
}

func TestCreateProductInvalidPrice(t *testing.T) {
	repo := &mocks.ProductRepository{CreateFn: func(context.Context, *entities.Product) error { return nil }}
	uc := productuc.NewUseCase(repo)
	_, err := uc.CreateProduct(context.Background(), dto.CreateProductInput{Name: "Menu", ProductType: "LUNCH", BasePrice: decimal.Zero})
	if !errors.Is(err, domainerrors.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestListProductsFiltered(t *testing.T) {
	pt1, _ := valueobjects.NewProductType("LUNCH")
	pt2, _ := valueobjects.NewProductType("JUICE")
	repo := &mocks.ProductRepository{
		ListActiveFn: func(context.Context) ([]entities.Product, error) {
			return []entities.Product{{ID: uuid.New(), ProductType: pt1, IsActive: true}, {ID: uuid.New(), ProductType: pt2, IsActive: true}}, nil
		},
	}
	uc := productuc.NewUseCase(repo)
	items, err := uc.ListProducts(context.Background(), dto.ListProductsInput{ProductType: "LUNCH"})
	if err != nil || len(items) != 1 {
		t.Fatalf("unexpected result: %v len=%d", err, len(items))
	}
}

func TestListPublicProductsFiltered(t *testing.T) {
	pt1, _ := valueobjects.NewProductType("LUNCH")
	pt2, _ := valueobjects.NewProductType("JUICE")
	repo := &mocks.ProductRepository{
		ListPublicFn: func(context.Context) ([]entities.Product, error) {
			return []entities.Product{
				{ID: uuid.New(), Name: "Menu", ProductType: pt1, BasePrice: decimal.RequireFromString("10000"), ImageURL: "https://cdn.example.com/menu.jpg", IsActive: true, IsPublic: true},
				{ID: uuid.New(), Name: "Jugo", ProductType: pt2, BasePrice: decimal.RequireFromString("8000"), IsActive: true, IsPublic: true},
			}, nil
		},
	}
	uc := productuc.NewUseCase(repo)
	items, err := uc.ListPublicProducts(context.Background(), dto.ListProductsInput{ProductType: "LUNCH"})
	if err != nil || len(items) != 1 {
		t.Fatalf("unexpected result: %v len=%d", err, len(items))
	}
	if items[0].ImageURL != "https://cdn.example.com/menu.jpg" {
		t.Fatalf("unexpected public output: %+v", items[0])
	}
}

func TestListPublicCategories(t *testing.T) {
	lunch, _ := valueobjects.NewProductType("LUNCH")
	juice, _ := valueobjects.NewProductType("JUICE")
	repo := &mocks.ProductRepository{
		ListPublicFn: func(context.Context) ([]entities.Product, error) {
			return []entities.Product{
				{ID: uuid.New(), ProductType: lunch, IsActive: true, IsPublic: true},
				{ID: uuid.New(), ProductType: juice, IsActive: true, IsPublic: true},
				{ID: uuid.New(), ProductType: lunch, IsActive: true, IsPublic: true},
			}, nil
		},
	}
	uc := productuc.NewUseCase(repo)
	categories, err := uc.ListPublicCategories(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(categories) != 2 || categories[0] != "JUICE" || categories[1] != "LUNCH" {
		t.Fatalf("unexpected categories: %+v", categories)
	}
}

func TestUpdateProductSuccess(t *testing.T) {
	pt, _ := valueobjects.NewProductType("LUNCH")
	repo := &mocks.ProductRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Product, error) {
			return &entities.Product{ID: uuid.New(), ProductType: pt, IsActive: true}, nil
		},
		UpdateFn: func(context.Context, *entities.Product) error { return nil },
	}
	uc := productuc.NewUseCase(repo)
	_, err := uc.UpdateProduct(context.Background(), dto.UpdateProductInput{
		ID: uuid.New(), Name: "Nuevo", ProductType: "CAKE", BasePrice: decimal.RequireFromString("15000"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeactivateProductInactive(t *testing.T) {
	repo := &mocks.ProductRepository{
		FindByIDFn: func(context.Context, uuid.UUID) (*entities.Product, error) {
			pt, _ := valueobjects.NewProductType("LUNCH")
			return &entities.Product{ID: uuid.New(), ProductType: pt, IsActive: false, UpdatedAt: time.Now()}, nil
		},
		DeactivateFn: func(context.Context, uuid.UUID, time.Time) error { return nil },
	}
	uc := productuc.NewUseCase(repo)
	err := uc.DeactivateProduct(context.Background(), uuid.New())
	if !errors.Is(err, domainerrors.ErrInactive) {
		t.Fatalf("expected inactive error, got %v", err)
	}
}
