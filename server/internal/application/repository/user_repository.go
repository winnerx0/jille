package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/domain"
)

type UserRepository interface {
	FindById(ctx context.Context, userID uuid.UUID) (domain.User, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	FindByEmail(ctx context.Context, email string) (domain.User, error)

	Save(ctx context.Context, user *domain.User) error
}
