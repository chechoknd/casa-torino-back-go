package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/casatorino/backend/internal/domain/valueobjects"
)

type RegisterUserInput struct {
	Email    string
	Username string
	FullName string
	Password string
}

type LoginInput struct {
	EmailOrUsername string
	Password        string
}

type AuthUserOutput struct {
	ID        uuid.UUID             `json:"id"`
	Email     string                `json:"email"`
	Username  string                `json:"username"`
	FullName  string                `json:"full_name"`
	Role      valueobjects.UserRole `json:"role"`
	CreatedAt time.Time             `json:"created_at"`
}

type AuthTokenOutput struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token,omitempty"`
	TokenType    string         `json:"token_type"`
	ExpiresAt    time.Time      `json:"expires_at"`
	User         AuthUserOutput `json:"user"`
}

type RefreshTokenInput struct {
	RefreshToken string
}

type LogoutInput struct {
	RefreshToken string
}
