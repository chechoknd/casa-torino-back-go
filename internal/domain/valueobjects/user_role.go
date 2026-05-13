package valueobjects

import (
	"fmt"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleCustomer UserRole = "CUSTOMER"
)

func NewUserRole(value string) (UserRole, error) {
	role := UserRole(value)
	switch role {
	case UserRoleAdmin, UserRoleCustomer:
		return role, nil
	default:
		return "", fmt.Errorf("%w: user role %q", domainerrors.ErrInvalidInput, value)
	}
}
