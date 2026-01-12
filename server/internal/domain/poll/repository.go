package poll

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {

	FindUserPollCount(ctx context.Context, userID uuid.UUID) (int, error)
}