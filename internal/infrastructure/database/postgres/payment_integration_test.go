package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/application/dto"
	paymentuc "github.com/casatorino/backend/internal/application/usecases/payment"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
	"github.com/casatorino/backend/internal/infrastructure/database/postgres"
)

type paymentIntegrationHarness struct {
	customerRepo *postgres.CustomerRepository
	productRepo  *postgres.ProductRepository
	orderRepo    *postgres.OrderRepository
	paymentRepo  *postgres.PaymentRepository
	useCase      *paymentuc.UseCase
}

func TestPaymentIntegrationCreatePersistsAndEnrichesOutput(t *testing.T) {
	h := newPaymentIntegrationHarness(t)
	ctx := context.Background()
	order := h.createOrderWithItem(t, ctx)

	out, err := h.useCase.CreatePayment(ctx, dto.CreatePaymentInput{
		OrderID: order.ID,
		Amount:  decimal.RequireFromString("35000"),
		Method:  "CASH",
		Status:  "PENDING",
	})
	if err != nil {
		t.Fatalf("create payment: %v", err)
	}

	if out.ID == uuid.Nil {
		t.Fatal("expected generated payment id")
	}
	if out.OrderID != order.ID || out.OrderNumber != order.OrderNumber || out.OrderLabel == "" {
		t.Fatalf("unexpected order data in output: %+v", out)
	}
	if !out.Amount.Equal(decimal.RequireFromString("35000")) || out.Method != "CASH" || out.Status != "PENDING" {
		t.Fatalf("unexpected payment output: %+v", out)
	}
	if len(out.Products) != 1 || out.Products[0].ProductName != "Integration Product" || out.Products[0].Quantity != 2 {
		t.Fatalf("unexpected payment products: %+v", out.Products)
	}

	persisted, err := h.paymentRepo.FindByID(ctx, out.ID)
	if err != nil {
		t.Fatalf("find persisted payment: %v", err)
	}
	if persisted.OrderID != order.ID || !persisted.Amount.Equal(decimal.RequireFromString("35000")) {
		t.Fatalf("unexpected persisted payment: %+v", persisted)
	}
}

func TestPaymentIntegrationListAndListByOrder(t *testing.T) {
	h := newPaymentIntegrationHarness(t)
	ctx := context.Background()
	firstOrder := h.createOrderWithItem(t, ctx)
	secondOrder := h.createOrderWithItem(t, ctx)

	firstPayment := h.createPayment(t, ctx, firstOrder.ID, "11000", "PAID")
	secondPayment := h.createPayment(t, ctx, firstOrder.ID, "12000", "PENDING")
	h.createPayment(t, ctx, secondOrder.ID, "13000", "PAID")

	byOrder, err := h.useCase.GetPaymentsByOrder(ctx, firstOrder.ID)
	if err != nil {
		t.Fatalf("get payments by order: %v", err)
	}
	if len(byOrder) != 2 {
		t.Fatalf("payments by order len = %d, want 2: %+v", len(byOrder), byOrder)
	}
	assertPaymentIDs(t, byOrder, firstPayment.ID, secondPayment.ID)

	all, err := h.useCase.ListPayments(ctx)
	if err != nil {
		t.Fatalf("list payments: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("all payments len = %d, want 3: %+v", len(all), all)
	}
}

func TestPaymentIntegrationUpdateStatusPersistsFinalState(t *testing.T) {
	h := newPaymentIntegrationHarness(t)
	ctx := context.Background()
	order := h.createOrderWithItem(t, ctx)
	payment := h.createPayment(t, ctx, order.ID, "35000", "PENDING")

	out, err := h.useCase.UpdatePaymentStatus(ctx, dto.UpdatePaymentStatusInput{
		PaymentID: payment.ID,
		Status:    "PAID",
	})
	if err != nil {
		t.Fatalf("update payment status: %v", err)
	}
	if out.Status != "PAID" {
		t.Fatalf("output status = %s, want PAID", out.Status)
	}

	persisted, err := h.paymentRepo.FindByID(ctx, payment.ID)
	if err != nil {
		t.Fatalf("find updated payment: %v", err)
	}
	if persisted.Status != valueobjects.PaymentStatusPaid {
		t.Fatalf("persisted status = %s, want PAID", persisted.Status)
	}
	if !persisted.UpdatedAt.After(payment.UpdatedAt) {
		t.Fatalf("expected updated_at to move forward: before=%s after=%s", payment.UpdatedAt, persisted.UpdatedAt)
	}
}

func TestPaymentIntegrationValidationFailuresDoNotPersist(t *testing.T) {
	h := newPaymentIntegrationHarness(t)
	ctx := context.Background()
	order := h.createOrderWithItem(t, ctx)

	tests := []struct {
		name  string
		input dto.CreatePaymentInput
	}{
		{
			name:  "zero amount",
			input: dto.CreatePaymentInput{OrderID: order.ID, Amount: decimal.Zero, Method: "CASH", Status: "PENDING"},
		},
		{
			name:  "negative amount",
			input: dto.CreatePaymentInput{OrderID: order.ID, Amount: decimal.RequireFromString("-1"), Method: "CASH", Status: "PENDING"},
		},
		{
			name:  "invalid method",
			input: dto.CreatePaymentInput{OrderID: order.ID, Amount: decimal.RequireFromString("1000"), Method: "BITCOIN", Status: "PENDING"},
		},
		{
			name:  "invalid status",
			input: dto.CreatePaymentInput{OrderID: order.ID, Amount: decimal.RequireFromString("1000"), Method: "CASH", Status: "UNKNOWN"},
		},
		{
			name:  "missing order",
			input: dto.CreatePaymentInput{OrderID: uuid.New(), Amount: decimal.RequireFromString("1000"), Method: "CASH", Status: "PENDING"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := h.useCase.CreatePayment(ctx, tc.input); err == nil {
				t.Fatal("expected error")
			}
		})
	}

	payments, err := h.paymentRepo.ListByOrderID(ctx, order.ID)
	if err != nil {
		t.Fatalf("list payments after validation failures: %v", err)
	}
	if len(payments) != 0 {
		t.Fatalf("expected no persisted payments, got %+v", payments)
	}
}

func TestPaymentIntegrationRepositoryRejectsOrphanPayment(t *testing.T) {
	h := newPaymentIntegrationHarness(t)
	ctx := context.Background()
	method, _ := valueobjects.NewPaymentMethod("CASH")
	status, _ := valueobjects.NewPaymentStatus("PENDING")
	now := time.Now().UTC()
	payment := &entities.Payment{
		ID:        uuid.New(),
		OrderID:   uuid.New(),
		Amount:    decimal.RequireFromString("1000"),
		Method:    method,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := h.paymentRepo.Create(ctx, payment)
	if !errors.Is(err, domainerrors.ErrInvalidInput) {
		t.Fatalf("error = %v, want ErrInvalidInput", err)
	}
}

func TestPaymentIntegrationConcurrentCreatesPersistAllPayments(t *testing.T) {
	h := newPaymentIntegrationHarness(t)
	ctx := context.Background()
	order := h.createOrderWithItem(t, ctx)
	const workers = 8

	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for index := 0; index < workers; index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := h.useCase.CreatePayment(ctx, dto.CreatePaymentInput{
				OrderID: order.ID,
				Amount:  decimal.NewFromInt(int64(1000 + index)),
				Method:  "CASH",
				Status:  "PENDING",
			})
			errs <- err
		}(index)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent create payment: %v", err)
		}
	}

	payments, err := h.paymentRepo.ListByOrderID(ctx, order.ID)
	if err != nil {
		t.Fatalf("list payments after concurrent creates: %v", err)
	}
	if len(payments) != workers {
		t.Fatalf("payments len = %d, want %d", len(payments), workers)
	}
}

func newPaymentIntegrationHarness(t *testing.T) *paymentIntegrationHarness {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("set TEST_DATABASE_URL or DATABASE_URL to run PostgreSQL integration tests")
	}

	ctx := context.Background()
	adminPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect admin pool: %v", err)
	}
	t.Cleanup(adminPool.Close)

	schema := "test_payment_" + strings.ReplaceAll(uuid.NewString(), "-", "_")
	if _, err := adminPool.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatalf("create test schema: %v", err)
	}
	t.Cleanup(func() {
		_, _ = adminPool.Exec(context.Background(), "DROP SCHEMA IF EXISTS "+schema+" CASCADE")
	})

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatalf("parse database url: %v", err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("connect test pool: %v", err)
	}
	t.Cleanup(pool.Close)

	applyPaymentTestMigrations(t, ctx, pool)

	customerRepo := postgres.NewCustomerRepository(pool)
	productRepo := postgres.NewProductRepository(pool)
	orderRepo := postgres.NewOrderRepository(pool)
	paymentRepo := postgres.NewPaymentRepository(pool)

	return &paymentIntegrationHarness{
		customerRepo: customerRepo,
		productRepo:  productRepo,
		orderRepo:    orderRepo,
		paymentRepo:  paymentRepo,
		useCase:      paymentuc.NewUseCase(paymentRepo, orderRepo, productRepo),
	}
}

func applyPaymentTestMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	migrationsDir := filepath.Join("..", "..", "..", "..", "migrations")
	for _, name := range []string{
		"000001_create_core_tables.up.sql",
		"000002_add_order_number.up.sql",
	} {
		content, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			t.Fatalf("read migration %s: %v", name, err)
		}
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			t.Fatalf("apply migration %s: %v", name, err)
		}
	}
}

func (h *paymentIntegrationHarness) createOrderWithItem(t *testing.T, ctx context.Context) entities.Order {
	t.Helper()

	customerType, _ := valueobjects.NewCustomerType("PERSON")
	now := time.Now().UTC()
	customer := &entities.Customer{
		ID:           uuid.New(),
		FullName:     "Integration Customer",
		Phone:        "3000000000",
		Email:        fmt.Sprintf("%s@example.com", uuid.NewString()),
		CustomerType: customerType,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}
	if err := h.customerRepo.Create(ctx, customer); err != nil {
		t.Fatalf("create customer fixture: %v", err)
	}

	productType, _ := valueobjects.NewProductType("LUNCH")
	product := &entities.Product{
		ID:          uuid.New(),
		Name:        "Integration Product",
		Description: "Payment integration product",
		ProductType: productType,
		BasePrice:   decimal.RequireFromString("18000"),
		CostPrice:   decimal.RequireFromString("9000"),
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := h.productRepo.Create(ctx, product); err != nil {
		t.Fatalf("create product fixture: %v", err)
	}

	order := &entities.Order{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		Status:     valueobjects.OrderStatusPending,
		Items: []entities.OrderItem{{
			ID:        uuid.New(),
			ProductID: product.ID,
			Quantity:  2,
			UnitPrice: product.BasePrice,
		}},
		Discount:  decimal.RequireFromString("1000"),
		CreatedAt: now,
		UpdatedAt: now,
	}
	order.CalculateTotal()
	if err := h.orderRepo.Create(ctx, order); err != nil {
		t.Fatalf("create order fixture: %v", err)
	}
	if order.OrderNumber == 0 {
		t.Fatal("expected database-generated order number")
	}
	return *order
}

func (h *paymentIntegrationHarness) createPayment(t *testing.T, ctx context.Context, orderID uuid.UUID, amount string, status string) entities.Payment {
	t.Helper()

	out, err := h.useCase.CreatePayment(ctx, dto.CreatePaymentInput{
		OrderID: orderID,
		Amount:  decimal.RequireFromString(amount),
		Method:  "CASH",
		Status:  status,
	})
	if err != nil {
		t.Fatalf("create payment fixture: %v", err)
	}

	payment, err := h.paymentRepo.FindByID(ctx, out.ID)
	if err != nil {
		t.Fatalf("find payment fixture: %v", err)
	}
	return *payment
}

func assertPaymentIDs(t *testing.T, payments []dto.PaymentOutput, ids ...uuid.UUID) {
	t.Helper()

	seen := make(map[uuid.UUID]struct{}, len(payments))
	for _, payment := range payments {
		seen[payment.ID] = struct{}{}
	}
	for _, id := range ids {
		if _, ok := seen[id]; !ok {
			t.Fatalf("missing payment id %s in %+v", id, payments)
		}
	}
}
