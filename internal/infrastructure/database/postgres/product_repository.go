package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casatorino/backend/internal/domain/entities"
	domainrepositories "github.com/casatorino/backend/internal/domain/repositories"
	sqlcdb "github.com/casatorino/backend/internal/infrastructure/database/sqlc"
)

var _ domainrepositories.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	conn    *pgxpool.Pool
	queries *sqlcdb.Queries
}

func NewProductRepository(conn *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		conn:    conn,
		queries: sqlcdb.New(conn),
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *entities.Product) error {
	row, err := r.queries.CreateProduct(ctx, sqlcdb.CreateProductParams{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ProductType: string(product.ProductType),
		BasePrice:   product.BasePrice,
		CostPrice:   product.CostPrice,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	})
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapProduct(row)
	if err != nil {
		return err
	}

	*product = mapped
	return nil
}

func (r *ProductRepository) Update(ctx context.Context, product *entities.Product) error {
	row, err := r.queries.UpdateProduct(ctx, sqlcdb.UpdateProductParams{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ProductType: string(product.ProductType),
		BasePrice:   product.BasePrice,
		CostPrice:   product.CostPrice,
		IsActive:    product.IsActive,
		UpdatedAt:   product.UpdatedAt,
	})
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapProduct(row)
	if err != nil {
		return err
	}

	*product = mapped
	return nil
}

func (r *ProductRepository) Deactivate(ctx context.Context, id uuid.UUID, updatedAt time.Time) error {
	rows, err := r.queries.DeactivateProduct(ctx, sqlcdb.DeactivateProductParams{
		ID:        id,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return mapError(err)
	}
	if rows == 0 {
		return domainrepositoriesErrNotFound()
	}
	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	row, err := r.queries.GetProductByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	mapped, err := mapProduct(row)
	if err != nil {
		return nil, err
	}

	return &mapped, nil
}

func (r *ProductRepository) ListActive(ctx context.Context) ([]entities.Product, error) {
	rows, err := r.queries.ListProducts(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	products := make([]entities.Product, 0, len(rows))
	for _, row := range rows {
		product, err := mapProduct(row)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
