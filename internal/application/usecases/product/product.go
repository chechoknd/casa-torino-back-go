package product

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
	products repositories.ProductRepository
}

func NewUseCase(products repositories.ProductRepository) *UseCase {
	return &UseCase{products: products}
}

func (uc *UseCase) CreateProduct(ctx context.Context, input dto.CreateProductInput) (dto.ProductOutput, error) {
	if strings.TrimSpace(input.Name) == "" || !input.BasePrice.IsPositive() {
		return dto.ProductOutput{}, domainerrors.ErrInvalidInput
	}

	productType, err := valueobjects.NewProductType(input.ProductType)
	if err != nil {
		return dto.ProductOutput{}, err
	}

	now := time.Now().UTC()
	product := &entities.Product{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		ProductType: productType,
		BasePrice:   input.BasePrice,
		CostPrice:   input.CostPrice,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.products.Create(ctx, product); err != nil {
		return dto.ProductOutput{}, err
	}

	return toProductOutput(*product), nil
}

func (uc *UseCase) GetProduct(ctx context.Context, id uuid.UUID) (dto.ProductOutput, error) {
	product, err := uc.products.FindByID(ctx, id)
	if err != nil {
		return dto.ProductOutput{}, err
	}
	if !product.IsActive {
		return dto.ProductOutput{}, domainerrors.ErrInactive
	}
	return toProductOutput(*product), nil
}

func (uc *UseCase) ListProducts(ctx context.Context, input dto.ListProductsInput) ([]dto.ProductOutput, error) {
	products, err := uc.products.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	var filter valueobjects.ProductType
	if strings.TrimSpace(input.ProductType) != "" {
		filter, err = valueobjects.NewProductType(input.ProductType)
		if err != nil {
			return nil, err
		}
	}

	output := make([]dto.ProductOutput, 0, len(products))
	for _, product := range products {
		if !product.IsActive {
			continue
		}
		if filter != "" && product.ProductType != filter {
			continue
		}
		output = append(output, toProductOutput(product))
	}

	return output, nil
}

func (uc *UseCase) UpdateProduct(ctx context.Context, input dto.UpdateProductInput) (dto.ProductOutput, error) {
	product, err := uc.products.FindByID(ctx, input.ID)
	if err != nil {
		return dto.ProductOutput{}, err
	}
	if !product.IsActive {
		return dto.ProductOutput{}, domainerrors.ErrInactive
	}
	if strings.TrimSpace(input.Name) == "" || !input.BasePrice.IsPositive() {
		return dto.ProductOutput{}, domainerrors.ErrInvalidInput
	}

	productType, err := valueobjects.NewProductType(input.ProductType)
	if err != nil {
		return dto.ProductOutput{}, err
	}

	product.Name = strings.TrimSpace(input.Name)
	product.Description = strings.TrimSpace(input.Description)
	product.ProductType = productType
	product.BasePrice = input.BasePrice
	product.CostPrice = input.CostPrice
	product.UpdatedAt = time.Now().UTC()

	if err := uc.products.Update(ctx, product); err != nil {
		return dto.ProductOutput{}, err
	}

	return toProductOutput(*product), nil
}

func (uc *UseCase) DeactivateProduct(ctx context.Context, id uuid.UUID) error {
	product, err := uc.products.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if !product.IsActive {
		return domainerrors.ErrInactive
	}

	return uc.products.Deactivate(ctx, id, time.Now().UTC())
}

func toProductOutput(product entities.Product) dto.ProductOutput {
	return dto.ProductOutput{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ProductType: string(product.ProductType),
		BasePrice:   product.BasePrice,
		CostPrice:   product.CostPrice,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
