package valueobjects

import (
	"fmt"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type Unit string

const (
	UnitG    Unit = "G"
	UnitML   Unit = "ML"
	UnitUnit Unit = "UNIT"
	UnitKG   Unit = "KG"
	UnitL    Unit = "L"
)

func NewUnit(value string) (Unit, error) {
	unit := Unit(value)
	switch unit {
	case UnitG, UnitML, UnitUnit, UnitKG, UnitL:
		return unit, nil
	default:
		return "", fmt.Errorf("%w: unit %q", domainerrors.ErrInvalidInput, value)
	}
}
