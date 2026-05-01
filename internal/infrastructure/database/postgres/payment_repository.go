package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casatorino/backend/internal/domain/entities"
	domainrepositories "github.com/casatorino/backend/internal/domain/repositories"
	sqlcdb "github.com/casatorino/backend/internal/infrastructure/database/sqlc"
)

var _ domainrepositories.PaymentRepository = (*PaymentRepository)(nil)

type PaymentRepository struct {
	conn    *pgxpool.Pool
	queries *sqlcdb.Queries
}

func NewPaymentRepository(conn *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{
		conn:    conn,
		queries: sqlcdb.New(conn),
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *entities.Payment) error {
	row, err := r.queries.CreatePayment(ctx, sqlcdb.CreatePaymentParams{
		ID:        payment.ID,
		OrderID:   payment.OrderID,
		Amount:    payment.Amount,
		Method:    string(payment.Method),
		Status:    string(payment.Status),
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	})
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapPayment(row)
	if err != nil {
		return err
	}

	*payment = mapped
	return nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *entities.Payment) error {
	row, err := r.queries.UpdatePayment(ctx, sqlcdb.UpdatePaymentParams{
		ID:        payment.ID,
		OrderID:   payment.OrderID,
		Amount:    payment.Amount,
		Method:    string(payment.Method),
		Status:    string(payment.Status),
		UpdatedAt: payment.UpdatedAt,
	})
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapPayment(row)
	if err != nil {
		return err
	}

	*payment = mapped
	return nil
}

func (r *PaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	row, err := r.queries.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	mapped, err := mapPayment(row)
	if err != nil {
		return nil, err
	}

	return &mapped, nil
}

func (r *PaymentRepository) ListByOrderID(ctx context.Context, orderID uuid.UUID) ([]entities.Payment, error) {
	rows, err := r.queries.GetPaymentsByOrderID(ctx, orderID)
	if err != nil {
		return nil, mapError(err)
	}

	payments := make([]entities.Payment, 0, len(rows))
	for _, row := range rows {
		payment, err := mapPayment(row)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) List(ctx context.Context) ([]entities.Payment, error) {
	rows, err := r.queries.ListPayments(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	payments := make([]entities.Payment, 0, len(rows))
	for _, row := range rows {
		payment, err := mapPayment(row)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, nil
}
