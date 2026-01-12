package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/common/dto"
)

type Service interface {
	GetUserById(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	GetUserByEmail(ctx context.Context, email string) (*dto.UserAuthView, error)
}
