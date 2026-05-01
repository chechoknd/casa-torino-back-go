package valueobjects

import (
	"fmt"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type ProductType string

const (
	ProductTypeLunch        ProductType = "LUNCH"
	ProductTypeJuice        ProductType = "JUICE"
	ProductTypeCake         ProductType = "CAKE"
	ProductTypeEvent        ProductType = "EVENT"
	ProductTypePlan         ProductType = "PLAN"
	ProductTypeVacuumPacked ProductType = "VACUUM_PACKED"
)

func NewProductType(value string) (ProductType, error) {
	productType := ProductType(value)
	switch productType {
	case ProductTypeLunch, ProductTypeJuice, ProductTypeCake, ProductTypeEvent, ProductTypePlan, ProductTypeVacuumPacked:
		return productType, nil
	default:
		return "", fmt.Errorf("%w: product type %q", domainerrors.ErrInvalidInput, value)
	}
}
