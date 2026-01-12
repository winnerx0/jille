package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`

	PollCount int `json:"poll_count"`

	ProfilePicture string `json:"profile_picture,omitempty"`
}

type CreateUserRequest struct {
	Username string `validate:"required,min=5,max=10"`

	Email string `validate:"required,email"`

	Password string `validate:"required,min=8,max=16"`
}

type LoginUserRequest struct {
	Email string `validate:"required,email"`

	Password string `validate:"required"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	Message    string `json:"message"`
	AuthTokens `json:"data"`
}

type UserAuthView struct {
	ID       uuid.UUID
	Email    string
	Password string
}
