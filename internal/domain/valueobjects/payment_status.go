package valueobjects

import (
	"fmt"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusPartial  PaymentStatus = "PARTIAL"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

func NewPaymentStatus(value string) (PaymentStatus, error) {
	status := PaymentStatus(value)
	switch status {
	case PaymentStatusPending, PaymentStatusPaid, PaymentStatusPartial, PaymentStatusFailed, PaymentStatusRefunded:
		return status, nil
	default:
		return "", fmt.Errorf("%w: payment status %q", domainerrors.ErrInvalidInput, value)
	}
}
