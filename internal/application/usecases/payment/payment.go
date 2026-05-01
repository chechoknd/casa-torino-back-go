package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/casatorino/backend/internal/application/dto"
	"github.com/casatorino/backend/internal/domain/entities"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/repositories"
	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type UseCase struct {
	payments repositories.PaymentRepository
	orders   repositories.OrderRepository
}

func NewUseCase(payments repositories.PaymentRepository, orders repositories.OrderRepository) *UseCase {
	return &UseCase{
		payments: payments,
		orders:   orders,
	}
}

func (uc *UseCase) CreatePayment(ctx context.Context, input dto.CreatePaymentInput) (dto.PaymentOutput, error) {
	if input.Amount.LessThanOrEqual(decimal.Zero) {
		return dto.PaymentOutput{}, domainerrors.ErrInvalidInput
	}

	if _, err := uc.orders.FindByID(ctx, input.OrderID); err != nil {
		return dto.PaymentOutput{}, err
	}

	method, err := valueobjects.NewPaymentMethod(input.Method)
	if err != nil {
		return dto.PaymentOutput{}, err
	}
	status, err := valueobjects.NewPaymentStatus(input.Status)
	if err != nil {
		return dto.PaymentOutput{}, err
	}

	now := time.Now().UTC()
	payment := &entities.Payment{
		ID:        uuid.New(),
		OrderID:   input.OrderID,
		Amount:    input.Amount,
		Method:    method,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.payments.Create(ctx, payment); err != nil {
		return dto.PaymentOutput{}, err
	}

	return toPaymentOutput(*payment), nil
}

func (uc *UseCase) GetPaymentsByOrder(ctx context.Context, orderID uuid.UUID) ([]dto.PaymentOutput, error) {
	if _, err := uc.orders.FindByID(ctx, orderID); err != nil {
		return nil, err
	}

	payments, err := uc.payments.ListByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	output := make([]dto.PaymentOutput, 0, len(payments))
	for _, payment := range payments {
		output = append(output, toPaymentOutput(payment))
	}

	return output, nil
}

func (uc *UseCase) UpdatePaymentStatus(ctx context.Context, input dto.UpdatePaymentStatusInput) (dto.PaymentOutput, error) {
	payment, err := uc.payments.FindByID(ctx, input.PaymentID)
	if err != nil {
		return dto.PaymentOutput{}, err
	}

	status, err := valueobjects.NewPaymentStatus(input.Status)
	if err != nil {
		return dto.PaymentOutput{}, err
	}

	payment.Status = status
	payment.UpdatedAt = time.Now().UTC()

	if err := uc.payments.Update(ctx, payment); err != nil {
		return dto.PaymentOutput{}, err
	}

	return toPaymentOutput(*payment), nil
}

func toPaymentOutput(payment entities.Payment) dto.PaymentOutput {
	return dto.PaymentOutput{
		ID:        payment.ID,
		OrderID:   payment.OrderID,
		Amount:    payment.Amount,
		Method:    string(payment.Method),
		Status:    string(payment.Status),
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}
}
