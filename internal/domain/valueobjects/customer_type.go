package valueobjects

import (
	"fmt"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type CustomerType string

const (
	CustomerTypePerson  CustomerType = "PERSON"
	CustomerTypeCompany CustomerType = "COMPANY"
)

func NewCustomerType(value string) (CustomerType, error) {
	customerType := CustomerType(value)
	switch customerType {
	case CustomerTypePerson, CustomerTypeCompany:
		return customerType, nil
	default:
		return "", fmt.Errorf("%w: customer type %q", domainerrors.ErrInvalidInput, value)
	}
}
