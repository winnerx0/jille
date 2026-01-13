package repository

import (
	"context"

	"github.com/winnerx0/jille/internal/domain"
)

type OptionRepository interface {
	Save(ctx context.Context, option *[]domain.Option) error
}
