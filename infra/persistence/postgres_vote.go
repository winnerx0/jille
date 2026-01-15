package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/domain"
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

	return err
}