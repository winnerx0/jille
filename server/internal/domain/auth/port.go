package auth

import (
	"context"

	"github.com/winnerx0/jille/internal/common/dto"
)

type Service interface {
	Login(ctx context.Context, loginRequest dto.LoginUserRequest) (*dto.AuthResponse, error)
}
