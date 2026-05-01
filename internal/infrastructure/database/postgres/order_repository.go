package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/casatorino/backend/internal/domain/entities"
	domainrepositories "github.com/casatorino/backend/internal/domain/repositories"
	sqlcdb "github.com/casatorino/backend/internal/infrastructure/database/sqlc"
)

var _ domainrepositories.OrderRepository = (*OrderRepository)(nil)

type OrderRepository struct {
	conn    *pgxpool.Pool
	queries *sqlcdb.Queries
}

func NewOrderRepository(conn *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		conn:    conn,
		queries: sqlcdb.New(conn),
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *entities.Order) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)
	row, err := qtx.CreateOrder(ctx, sqlcdb.CreateOrderParams{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     string(order.Status),
		Subtotal:   order.Subtotal,
		Discount:   order.Discount,
		Total:      order.Total,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	})
	if err != nil {
		return mapError(err)
	}

	for _, item := range order.Items {
		if _, err := qtx.CreateOrderItem(ctx, sqlcdb.CreateOrderItemParams{
			ID:        item.ID,
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  int32(item.Quantity),
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		}); err != nil {
			return mapError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	items, err := r.queries.GetOrderItemsByOrderID(ctx, row.ID)
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapOrder(row, items)
	if err != nil {
		return err
	}

	*order = mapped
	return nil
}

func (r *OrderRepository) Update(ctx context.Context, order *entities.Order) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)
	row, err := qtx.UpdateOrder(ctx, sqlcdb.UpdateOrderParams{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     string(order.Status),
		Subtotal:   order.Subtotal,
		Discount:   order.Discount,
		Total:      order.Total,
		UpdatedAt:  order.UpdatedAt,
	})
	if err != nil {
		return mapError(err)
	}

	for _, item := range order.Items {
		if item.ID == uuid.Nil {
			continue
		}

		if _, err := qtx.UpdateOrderItem(ctx, sqlcdb.UpdateOrderItemParams{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  int32(item.Quantity),
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		}); err != nil {
			return mapError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	items, err := r.queries.GetOrderItemsByOrderID(ctx, row.ID)
	if err != nil {
		return mapError(err)
	}

	mapped, err := mapOrder(row, items)
	if err != nil {
		return err
	}

	*order = mapped
	return nil
}

func (r *OrderRepository) AddItem(ctx context.Context, orderID uuid.UUID, item *entities.OrderItem) error {
	row, err := r.queries.CreateOrderItem(ctx, sqlcdb.CreateOrderItemParams{
		ID:        item.ID,
		OrderID:   orderID,
		ProductID: item.ProductID,
		Quantity:  int32(item.Quantity),
		UnitPrice: item.UnitPrice,
		Subtotal:  item.Subtotal,
	})
	if err != nil {
		return mapError(err)
	}

	item.ID = row.ID
	item.OrderID = row.OrderID
	item.ProductID = row.ProductID
	item.Quantity = int(row.Quantity)
	item.UnitPrice = row.UnitPrice
	item.Subtotal = row.Subtotal

	return nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	row, err := r.queries.GetOrderByID(ctx, id)
	if err != nil {
		return nil, mapError(err)
	}

	items, err := r.queries.GetOrderItemsByOrderID(ctx, row.ID)
	if err != nil {
		return nil, mapError(err)
	}

	mapped, err := mapOrder(row, items)
	if err != nil {
		return nil, err
	}

	return &mapped, nil
}

func (r *OrderRepository) List(ctx context.Context) ([]entities.Order, error) {
	rows, err := r.queries.ListOrders(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	orders := make([]entities.Order, 0, len(rows))
	for _, row := range rows {
		items, err := r.queries.GetOrderItemsByOrderID(ctx, row.ID)
		if err != nil {
			return nil, mapError(err)
		}

		order, err := mapOrder(row, items)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]entities.Order, error) {
	rows, err := r.queries.ListOrdersByCustomerID(ctx, customerID)
	if err != nil {
		return nil, mapError(err)
	}

	orders := make([]entities.Order, 0, len(rows))
	for _, row := range rows {
		items, err := r.queries.GetOrderItemsByOrderID(ctx, row.ID)
		if err != nil {
			return nil, mapError(err)
		}

		order, err := mapOrder(row, items)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
