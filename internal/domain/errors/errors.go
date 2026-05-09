package errors

import stderrors "errors"

var (
	ErrNotFound           = stderrors.New("domain: not found")
	ErrInvalidInput       = stderrors.New("domain: invalid input")
	ErrInvalidStatus      = stderrors.New("domain: invalid status")
	ErrInactive           = stderrors.New("domain: inactive entity")
	ErrDuplicateEmail     = stderrors.New("domain: duplicate email")
	ErrDuplicateUsername  = stderrors.New("domain: duplicate username")
	ErrInvalidCredentials = stderrors.New("domain: invalid credentials")
	ErrUnauthorized       = stderrors.New("domain: unauthorized")
	ErrRequestTooLarge    = stderrors.New("domain: request too large")
)
