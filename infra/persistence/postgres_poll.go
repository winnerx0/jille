package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/domain"
	"gorm.io/gorm"
)

type pollRepository struct {
	db *gorm.DB
}

func NewPollRepository(db *gorm.DB) repository.PollRepository {
	return &pollRepository{
		db: db,
	}
}

func (repo *pollRepository) FindUserPollCount(ctx context.Context, userID uuid.UUID) (int, error) {

	var pollCount int
	err := repo.db.
		Raw("SELECT COUNT(*) FROM polls WHERE user_id = ?", userID).
		Scan(&pollCount).Error
	if err != nil {
		return 0, err
	}

	return pollCount, nil
}

func (repo *pollRepository) Save(ctx context.Context, poll *domain.Poll) error {
	return gorm.G[domain.Poll](repo.db).Create(ctx, poll)
}

func (repo *pollRepository) FindPollByID(ctx context.Context, pollID uuid.UUID) (*domain.Poll, error) {
	var poll domain.Poll
	err := repo.db.First(&poll, pollID).Error
	if err != nil {
		return nil, err
	}
	return &poll, nil
}

func (repo *pollRepository) Delete(ctx context.Context, pollID uuid.UUID) error {
	rows, err := gorm.G[domain.Poll](repo.db).Where("id = ?", pollID).Delete(ctx)

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("poll not found")
	}

	return nil
}