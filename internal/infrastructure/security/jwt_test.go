package security

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
	"github.com/casatorino/backend/internal/domain/valueobjects"
)

func TestJWTManagerGenerateAndVerify(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	manager := NewJWTManagerWithClock("test-secret", 15*time.Minute, func() time.Time { return now })
	userID := uuid.New()

	token, expiresAt, err := manager.Generate(context.Background(), userID, "user@example.com", "demo", valueobjects.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Generate error = %v", err)
	}
	if !expiresAt.Equal(now.Add(15 * time.Minute)) {
		t.Fatalf("expiresAt = %s", expiresAt)
	}

	claims, err := manager.Verify(context.Background(), token)
	if err != nil {
		t.Fatalf("Verify error = %v", err)
	}
	if claims.UserID != userID || claims.Email != "user@example.com" || claims.Username != "demo" || claims.Role != valueobjects.UserRoleAdmin {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestJWTManagerRejectsExpiredToken(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	manager := NewJWTManagerWithClock("test-secret", time.Minute, func() time.Time { return now })
	token, _, err := manager.Generate(context.Background(), uuid.New(), "user@example.com", "demo", valueobjects.UserRoleCustomer)
	if err != nil {
		t.Fatalf("Generate error = %v", err)
	}

	expiredManager := NewJWTManagerWithClock("test-secret", time.Minute, func() time.Time {
		return now.Add(2 * time.Minute)
	})

	_, err = expiredManager.Verify(context.Background(), token)
	if !errors.Is(err, domainerrors.ErrUnauthorized) {
		t.Fatalf("error = %v, want ErrUnauthorized", err)
	}
}

func TestJWTManagerRejectsInvalidSignature(t *testing.T) {
	now := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)
	manager := NewJWTManagerWithClock("test-secret", time.Minute, func() time.Time { return now })
	token, _, err := manager.Generate(context.Background(), uuid.New(), "user@example.com", "demo", valueobjects.UserRoleCustomer)
	if err != nil {
		t.Fatalf("Generate error = %v", err)
	}

	otherManager := NewJWTManagerWithClock("other-secret", time.Minute, func() time.Time { return now })
	_, err = otherManager.Verify(context.Background(), token)
	if !errors.Is(err, domainerrors.ErrUnauthorized) {
		t.Fatalf("error = %v, want ErrUnauthorized", err)
	}
}
