package auth

import (
	"context"

	"github.com/winnerx0/jille/internal/common/dto"
)

type Service interface {
	Register(ctx context.Context, registerRequest dto.CreateUserRequest) (*dto.AuthResponse, error)

	Login(ctx context.Context, loginRequest dto.LoginUserRequest) (*dto.AuthResponse, error)
}
