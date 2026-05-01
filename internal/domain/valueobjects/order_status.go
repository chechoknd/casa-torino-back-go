package valueobjects

import (
	"fmt"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type OrderStatus string

const (
	OrderStatusPending       OrderStatus = "PENDING"
	OrderStatusConfirmed     OrderStatus = "CONFIRMED"
	OrderStatusInPreparation OrderStatus = "IN_PREPARATION"
	OrderStatusReady         OrderStatus = "READY"
	OrderStatusDelivered     OrderStatus = "DELIVERED"
	OrderStatusCancelled     OrderStatus = "CANCELLED"
)

var allowedOrderTransitions = map[OrderStatus]map[OrderStatus]struct{}{
	OrderStatusPending: {
		OrderStatusConfirmed: {},
		OrderStatusCancelled: {},
	},
	OrderStatusConfirmed: {
		OrderStatusInPreparation: {},
		OrderStatusCancelled:     {},
	},
	OrderStatusInPreparation: {
		OrderStatusReady:     {},
		OrderStatusCancelled: {},
	},
	OrderStatusReady: {
		OrderStatusDelivered: {},
	},
	OrderStatusDelivered: {},
	OrderStatusCancelled: {},
}

func NewOrderStatus(value string) (OrderStatus, error) {
	status := OrderStatus(value)
	if _, ok := allowedOrderTransitions[status]; !ok {
		return "", fmt.Errorf("%w: order status %q", domainerrors.ErrInvalidStatus, value)
	}

	return status, nil
}

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	if s == next {
		return true
	}

	transitions, ok := allowedOrderTransitions[s]
	if !ok {
		return false
	}

	_, allowed := transitions[next]
	return allowed
}
