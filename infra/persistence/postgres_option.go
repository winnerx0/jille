package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/domain"
	"gorm.io/gorm"
)

type optionRepository struct {
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) repository.OptionRepository {
	return &optionRepository{
		db: db,
	}
}

func (repo optionRepository) Save(ctx context.Context, options *[]domain.Option) error {

	return gorm.G[[]domain.Option](repo.db).Create(ctx, options)
}

func (v *optionRepository) FindOptionsByPollID(ctx context.Context, pollID uuid.UUID) (*[]domain.Option, error) {

	options, err := gorm.G[domain.Option](v.db).Preload("Votes", nil).Where("poll_id = ?", pollID).Find(ctx)

	if err != nil {
		return &[]domain.Option{}, err
	}

	return &options, nil
}
