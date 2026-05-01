package requests

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

func DecodeJSON(r *http.Request, target any) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return domainerrors.ErrInvalidInput
	}
	return nil
}

func ParseUUID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, domainerrors.ErrInvalidInput
	}
	return id, nil
}

func ParseDecimal(value string) (decimal.Decimal, error) {
	if value == "" {
		return decimal.Zero, nil
	}
	parsed, err := decimal.NewFromString(value)
	if err != nil {
		return decimal.Zero, domainerrors.ErrInvalidInput
	}
	return parsed, nil
}
