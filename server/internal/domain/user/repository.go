package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindById(ctx context.Context, userID uuid.UUID) (User, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	FindByEmail(ctx context.Context, email string) (User, error)
}
