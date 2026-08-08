package utils

import "time"

type Chirp struct {
	Body string `json:"body"`
}

type ValidateChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

// USERS
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	CreateUserRequest
}

type LoginRequest struct {
	CreateUserRequest
}

type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	UserResponse
}

// CHIRPS

type CreateChirpRequest struct {
	Body   string `json:"body"`
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type CreateChirpResponse struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// REFRESH TOKENS

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

// POLKA

type PolkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}
