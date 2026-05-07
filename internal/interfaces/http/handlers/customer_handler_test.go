package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/application/dto"
	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/interfaces/http/handlers"
	"github.com/casatorino/backend/internal/interfaces/http/routes"
)

type fakeCustomerUseCase struct {
	createFn     func(context.Context, dto.CreateCustomerInput) (dto.CustomerOutput, error)
	getFn        func(context.Context, uuid.UUID) (dto.CustomerOutput, error)
	listFn       func(context.Context) ([]dto.CustomerOutput, error)
	updateFn     func(context.Context, dto.UpdateCustomerInput) (dto.CustomerOutput, error)
	deactivateFn func(context.Context, uuid.UUID) error
}

func (f fakeCustomerUseCase) CreateCustomer(ctx context.Context, input dto.CreateCustomerInput) (dto.CustomerOutput, error) {
	return f.createFn(ctx, input)
}

func (f fakeCustomerUseCase) GetCustomer(ctx context.Context, id uuid.UUID) (dto.CustomerOutput, error) {
	return f.getFn(ctx, id)
}

func (f fakeCustomerUseCase) ListCustomers(ctx context.Context) ([]dto.CustomerOutput, error) {
	return f.listFn(ctx)
}

func (f fakeCustomerUseCase) UpdateCustomer(ctx context.Context, input dto.UpdateCustomerInput) (dto.CustomerOutput, error) {
	return f.updateFn(ctx, input)
}

func (f fakeCustomerUseCase) DeactivateCustomer(ctx context.Context, id uuid.UUID) error {
	return f.deactivateFn(ctx, id)
}

func TestCustomerEndpointsSuccess(t *testing.T) {
	customerID := uuid.New()
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	output := dto.CustomerOutput{
		ID:           customerID,
		FullName:     "Cliente Demo",
		Phone:        "3001234567",
		Email:        "cliente@example.com",
		CustomerType: "PERSON",
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}

	useCase := fakeCustomerUseCase{
		createFn: func(_ context.Context, input dto.CreateCustomerInput) (dto.CustomerOutput, error) {
			if input.FullName != "Cliente Demo" || input.Phone != "3001234567" || input.Email != "cliente@example.com" || input.CustomerType != "PERSON" {
				t.Fatalf("unexpected create input: %+v", input)
			}
			return output, nil
		},
		getFn: func(_ context.Context, id uuid.UUID) (dto.CustomerOutput, error) {
			if id != customerID {
				t.Fatalf("unexpected get id: %s", id)
			}
			return output, nil
		},
		listFn: func(context.Context) ([]dto.CustomerOutput, error) {
			return []dto.CustomerOutput{output}, nil
		},
		updateFn: func(_ context.Context, input dto.UpdateCustomerInput) (dto.CustomerOutput, error) {
			if input.ID != customerID || input.FullName != "Cliente Editado" || input.CustomerType != "COMPANY" {
				t.Fatalf("unexpected update input: %+v", input)
			}
			updated := output
			updated.FullName = input.FullName
			updated.CustomerType = input.CustomerType
			return updated, nil
		},
		deactivateFn: func(_ context.Context, id uuid.UUID) error {
			if id != customerID {
				t.Fatalf("unexpected delete id: %s", id)
			}
			return nil
		},
	}

	router := customerTestRouter(useCase)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{
			name:       "create",
			method:     http.MethodPost,
			path:       "/customers/",
			body:       `{"full_name":"Cliente Demo","phone":"3001234567","email":"cliente@example.com","customer_type":"PERSON"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "list",
			method:     http.MethodGet,
			path:       "/customers/",
			wantStatus: http.StatusOK,
		},
		{
			name:       "get",
			method:     http.MethodGet,
			path:       "/customers/" + customerID.String(),
			wantStatus: http.StatusOK,
		},
		{
			name:       "update",
			method:     http.MethodPut,
			path:       "/customers/" + customerID.String(),
			body:       `{"full_name":"Cliente Editado","phone":"3009990000","email":"editado@example.com","customer_type":"COMPANY"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "delete",
			method:     http.MethodDelete,
			path:       "/customers/" + customerID.String(),
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := performCustomerRequest(router, tc.method, tc.path, tc.body)
			if recorder.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body: %s", recorder.Code, tc.wantStatus, recorder.Body.String())
			}
			assertSuccessEnvelope(t, recorder.Body.Bytes())
		})
	}
}

func TestCustomerEndpointsErrors(t *testing.T) {
	customerID := uuid.New()
	useCase := fakeCustomerUseCase{
		createFn: func(context.Context, dto.CreateCustomerInput) (dto.CustomerOutput, error) {
			return dto.CustomerOutput{}, domainerrors.ErrDuplicateEmail
		},
		getFn: func(context.Context, uuid.UUID) (dto.CustomerOutput, error) {
			return dto.CustomerOutput{}, domainerrors.ErrNotFound
		},
		listFn: func(context.Context) ([]dto.CustomerOutput, error) {
			return nil, domainerrors.ErrInvalidInput
		},
		updateFn: func(context.Context, dto.UpdateCustomerInput) (dto.CustomerOutput, error) {
			return dto.CustomerOutput{}, domainerrors.ErrInactive
		},
		deactivateFn: func(context.Context, uuid.UUID) error {
			return domainerrors.ErrInactive
		},
	}

	router := customerTestRouter(useCase)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "create invalid json",
			method:     http.MethodPost,
			path:       "/customers/",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_INPUT",
		},
		{
			name:       "create duplicate email",
			method:     http.MethodPost,
			path:       "/customers/",
			body:       `{"full_name":"Cliente Demo","phone":"3001234567","email":"cliente@example.com","customer_type":"PERSON"}`,
			wantStatus: http.StatusConflict,
			wantCode:   "DUPLICATE_EMAIL",
		},
		{
			name:       "list usecase error",
			method:     http.MethodGet,
			path:       "/customers/",
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_INPUT",
		},
		{
			name:       "get invalid uuid",
			method:     http.MethodGet,
			path:       "/customers/not-a-uuid",
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_INPUT",
		},
		{
			name:       "get not found",
			method:     http.MethodGet,
			path:       "/customers/" + customerID.String(),
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
		},
		{
			name:       "update inactive",
			method:     http.MethodPut,
			path:       "/customers/" + customerID.String(),
			body:       `{"full_name":"Cliente Demo","phone":"3001234567","email":"cliente@example.com","customer_type":"PERSON"}`,
			wantStatus: http.StatusConflict,
			wantCode:   "INACTIVE",
		},
		{
			name:       "delete inactive",
			method:     http.MethodDelete,
			path:       "/customers/" + customerID.String(),
			wantStatus: http.StatusConflict,
			wantCode:   "INACTIVE",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := performCustomerRequest(router, tc.method, tc.path, tc.body)
			if recorder.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body: %s", recorder.Code, tc.wantStatus, recorder.Body.String())
			}
			assertErrorEnvelope(t, recorder.Body.Bytes(), tc.wantCode)
		})
	}
}

func customerTestRouter(useCase fakeCustomerUseCase) http.Handler {
	return routes.NewRouter(routes.Dependencies{
		Customers:   handlers.NewCustomerHandler(useCase),
		Products:    handlers.NewProductHandler(nil),
		Ingredients: handlers.NewIngredientHandler(nil),
		Recipes:     handlers.NewRecipeHandler(nil),
		Orders:      handlers.NewOrderHandler(nil),
		Payments:    handlers.NewPaymentHandler(nil),
	})
}

func performCustomerRequest(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	var requestBody *bytes.Reader
	if body == "" {
		requestBody = bytes.NewReader(nil)
	} else {
		requestBody = bytes.NewReader([]byte(body))
	}

	request := httptest.NewRequest(method, path, requestBody)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func assertSuccessEnvelope(t *testing.T, body []byte) {
	t.Helper()

	var envelope struct {
		Data    json.RawMessage `json:"data"`
		Message string          `json:"message"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("decode success response: %v", err)
	}
	if len(envelope.Data) == 0 {
		t.Fatalf("missing data field in response: %s", string(body))
	}
	if envelope.Message != "ok" {
		t.Fatalf("message = %q, want ok", envelope.Message)
	}
}

func assertErrorEnvelope(t *testing.T, body []byte, wantCode string) {
	t.Helper()

	var envelope struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if envelope.Code != wantCode {
		t.Fatalf("code = %q, want %q, body: %s", envelope.Code, wantCode, string(body))
	}
	if envelope.Error == "" {
		t.Fatalf("missing error message in response: %s", string(body))
	}
}
