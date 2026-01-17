package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/domain"
)

type OptionRepository interface {
	Save(ctx context.Context, option *[]domain.Option) error

	FindOptionsByPollID(ctx context.Context, pollID uuid.UUID) (*[]domain.Option, error)
}
