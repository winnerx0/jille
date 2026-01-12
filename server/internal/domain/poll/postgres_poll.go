package poll

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PollRepository struct {
	db *gorm.DB
}

func NewPollRepository(db *gorm.DB) Repository {
	return &PollRepository{
		db: db,
	}
}

func (repo *PollRepository) FindUserPollCount(ctx context.Context, userID uuid.UUID) (int, error) {

	var pollCount int
	err := repo.db.
		Raw("SELECT COUNT(*) FROM polls WHERE user_id = ?", userID).
		Scan(&pollCount).Error
	if err != nil {
		return 0, err
	}

	return pollCount, nil
}
