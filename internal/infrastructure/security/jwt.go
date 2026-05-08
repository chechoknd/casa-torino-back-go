package security

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	domainerrors "github.com/casatorino/backend/internal/domain/errors"
)

type TokenClaims struct {
	UserID    uuid.UUID
	Email     string
	Username  string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type JWTManager struct {
	secret    []byte
	expiresIn time.Duration
	now       func() time.Time
}

func NewJWTManager(secret string, expiresIn time.Duration) *JWTManager {
	return &JWTManager{
		secret:    []byte(secret),
		expiresIn: expiresIn,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

func NewJWTManagerWithClock(secret string, expiresIn time.Duration, now func() time.Time) *JWTManager {
	manager := NewJWTManager(secret, expiresIn)
	manager.now = now
	return manager
}

func (m *JWTManager) Generate(_ context.Context, userID uuid.UUID, email, username string) (string, time.Time, error) {
	now := m.now()
	expiresAt := now.Add(m.expiresIn)

	header := jwtHeader{
		Algorithm: "HS256",
		Type:      "JWT",
	}
	payload := jwtPayload{
		Subject:   userID.String(),
		Email:     email,
		Username:  username,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	unsigned, err := encodeSegments(header, payload)
	if err != nil {
		return "", time.Time{}, err
	}

	signature := sign(unsigned, m.secret)
	return unsigned + "." + signature, expiresAt, nil
}

func (m *JWTManager) Verify(_ context.Context, token string) (TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSignature := sign(unsigned, m.secret)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}

	var header jwtHeader
	if err := decodeSegment(parts[0], &header); err != nil {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}
	if header.Algorithm != "HS256" || header.Type != "JWT" {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}

	var payload jwtPayload
	if err := decodeSegment(parts[1], &payload); err != nil {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}
	if payload.Subject == "" || payload.Email == "" || payload.Username == "" || payload.ExpiresAt == 0 {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}
	if m.now().Unix() >= payload.ExpiresAt {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}

	userID, err := uuid.Parse(payload.Subject)
	if err != nil {
		return TokenClaims{}, domainerrors.ErrUnauthorized
	}

	return TokenClaims{
		UserID:    userID,
		Email:     payload.Email,
		Username:  payload.Username,
		IssuedAt:  time.Unix(payload.IssuedAt, 0).UTC(),
		ExpiresAt: time.Unix(payload.ExpiresAt, 0).UTC(),
	}, nil
}

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type jwtPayload struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func encodeSegments(header jwtHeader, payload jwtPayload) (string, error) {
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(headerBytes) + "." + base64.RawURLEncoding.EncodeToString(payloadBytes), nil
}

func decodeSegment(segment string, target any) error {
	decoded, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(decoded, target); err != nil {
		return err
	}

	return nil
}

func sign(unsigned string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
