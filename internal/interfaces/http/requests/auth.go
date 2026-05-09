package requests

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}
