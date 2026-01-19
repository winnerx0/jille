package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/domain"
)

type PollRepository interface {
	FindUserPollCount(ctx context.Context, userID uuid.UUID) (int, error)

	Save(ctx context.Context, poll *domain.Poll) error

	FindPollByID(ctx context.Context, pollID uuid.UUID) (*domain.Poll, error)

	Delete(ctx context.Context, pollID uuid.UUID) error

	FindAllPolls(ctx context.Context) ([]domain.Poll, error)
}
