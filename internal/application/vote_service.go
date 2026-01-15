package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/common/dto"
)

type VoteService interface {

	VotePoll(ctx context.Context, pollID uuid.UUID, optionID uuid.UUID) (*dto.VoteResponse, error)
}