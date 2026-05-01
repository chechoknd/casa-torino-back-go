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

var _ domainrepositories.CustomerRepository = (*CustomerRepository)(nil)

type CustomerRepository struct {
	conn    *pgxpool.Pool
	queries *sqlcdb.Queries
}

func NewCustomerRepository(conn *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{
		conn:    conn,
		queries: sqlcdb.New(conn),
	}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *entities.Customer) error {
	row, err := r.queries.CreateCustomer(ctx, sqlcdb.CreateCustomerParams{
		ID:           customer.ID,
		FullName:     customer.FullName,
		Phone:        customer.Phone,
		Email:        customer.Email,
		CustomerType: string(customer.CustomerType),
		CreatedAt:    customer.CreatedAt,
		UpdatedAt:    customer.UpdatedAt,
		IsActive:     customer.IsActive,
	})
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapCustomer(row)
	if err != nil {
		return err
	}

	*customer = mapped
	return nil
}

func (r *CustomerRepository) Update(ctx context.Context, customer *entities.Customer) error {
	row, err := r.queries.UpdateCustomer(ctx, sqlcdb.UpdateCustomerParams{
		ID:           customer.ID,
		FullName:     customer.FullName,
		Phone:        customer.Phone,
		Email:        customer.Email,
		CustomerType: string(customer.CustomerType),
		UpdatedAt:    customer.UpdatedAt,
		IsActive:     customer.IsActive,
	})
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapCustomer(row)
	if err != nil {
		return err
	}

	*customer = mapped
	return nil
}

func (r *CustomerRepository) Deactivate(ctx context.Context, id uuid.UUID, updatedAt time.Time) error {
	rows, err := r.queries.DeactivateCustomer(ctx, sqlcdb.DeactivateCustomerParams{
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

func (r *CustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error) {
	row, err := r.queries.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	mapped, err := mapCustomer(row)
	if err != nil {
		return nil, err
	}

	return &mapped, nil
}

func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	row, err := r.queries.GetCustomerByEmail(ctx, email)
	if err != nil {
		return nil, mapError(err)
	}

	mapped, err := mapCustomer(row)
	if err != nil {
		return nil, err
	}

	return &mapped, nil
}

func (r *CustomerRepository) List(ctx context.Context) ([]entities.Customer, error) {
	rows, err := r.queries.ListCustomers(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	customers := make([]entities.Customer, 0, len(rows))
	for _, row := range rows {
		customer, err := mapCustomer(row)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}
