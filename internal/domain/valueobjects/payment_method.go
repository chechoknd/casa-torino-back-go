package valueobjects

import (
	"fmt"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type PaymentMethod string

const (
	PaymentMethodCash      PaymentMethod = "CASH"
	PaymentMethodTransfer  PaymentMethod = "TRANSFER"
	PaymentMethodNequi     PaymentMethod = "NEQUI"
	PaymentMethodDaviplata PaymentMethod = "DAVIPLATA"
	PaymentMethodCard      PaymentMethod = "CARD"
	PaymentMethodOther     PaymentMethod = "OTHER"
)

func NewPaymentMethod(value string) (PaymentMethod, error) {
	method := PaymentMethod(value)
	switch method {
	case PaymentMethodCash, PaymentMethodTransfer, PaymentMethodNequi, PaymentMethodDaviplata, PaymentMethodCard, PaymentMethodOther:
		return method, nil
	default:
		return "", fmt.Errorf("%w: payment method %q", domainerrors.ErrInvalidInput, value)
	}
}
