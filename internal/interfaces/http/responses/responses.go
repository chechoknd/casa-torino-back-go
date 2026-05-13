package responses

import (
	"encoding/json"
	"errors"
	"net/http"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type SuccessResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(SuccessResponse{
		Data:    data,
		Message: "ok",
	})
}

func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := "internal server error"

	switch {
	case errors.Is(err, domainerrors.ErrNotFound):
		status = http.StatusNotFound
		code = "NOT_FOUND"
		message = "resource not found"
	case errors.Is(err, domainerrors.ErrInvalidInput):
		status = http.StatusBadRequest
		code = "INVALID_INPUT"
		message = "invalid input"
	case errors.Is(err, domainerrors.ErrInvalidStatus):
		status = http.StatusUnprocessableEntity
		code = "INVALID_STATUS"
		message = "invalid status"
	case errors.Is(err, domainerrors.ErrInactive):
		status = http.StatusConflict
		code = "INACTIVE"
		message = "resource is inactive"
	case errors.Is(err, domainerrors.ErrDuplicateEmail):
		status = http.StatusConflict
		code = "DUPLICATE_EMAIL"
		message = "email already exists"
	case errors.Is(err, domainerrors.ErrDuplicateUsername):
		status = http.StatusConflict
		code = "DUPLICATE_USERNAME"
		message = "username already exists"
	case errors.Is(err, domainerrors.ErrInvalidCredentials):
		status = http.StatusUnauthorized
		code = "INVALID_CREDENTIALS"
		message = "invalid credentials"
	case errors.Is(err, domainerrors.ErrUnauthorized):
		status = http.StatusUnauthorized
		code = "UNAUTHORIZED"
		message = "unauthorized"
	case errors.Is(err, domainerrors.ErrForbidden):
		status = http.StatusForbidden
		code = "FORBIDDEN"
		message = "forbidden"
	case errors.Is(err, domainerrors.ErrRequestTooLarge):
		status = http.StatusRequestEntityTooLarge
		code = "REQUEST_TOO_LARGE"
		message = "request entity too large"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}
