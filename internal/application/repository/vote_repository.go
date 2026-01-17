package repository

import (
	"context"

	"github.com/google/uuid"
)

type VoteRepository interface {
	Vote(ctx context.Context, pollID uuid.UUID, optionID uuid.UUID, userID uuid.UUID) error
}
