package persistence

import (
	"context"

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