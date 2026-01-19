package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/domain"
	"github.com/winnerx0/jille/internal/utils"
	"gorm.io/gorm"
)

type votereposutory struct {

	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) repository.VoteRepository {
	return &votereposutory{db: db}
}

func (v *votereposutory) Vote(ctx context.Context, pollID uuid.UUID, optionID uuid.UUID, userID uuid.UUID) error {


	err := gorm.G[domain.Vote](v.db).Create(ctx, &domain.Vote{
		PollID: pollID,
		UserID: userID,
		OptionID: optionID,
	})

	if err == gorm.ErrDuplicatedKey {
		return utils.VoteAlreadyExistsError
	}

	return err
}

func (v *votereposutory) ExistsByPollIDAndAndUserID(ctx context.Context, pollID uuid.UUID, userID uuid.UUID) (bool, error) {


	votes, err := gorm.G[domain.Vote](v.db).Where("poll_id = ? AND user_id = ?", pollID, userID).Find(ctx)

	if err == gorm.ErrDuplicatedKey {
		return false, err
	}

	return len(votes) > 0, nil
}