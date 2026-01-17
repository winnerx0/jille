package application

import (
	"context"

	"github.com/winnerx0/jille/internal/common/dto"
)

type VoteService interface {

	VotePoll(ctx context.Context, voteRequest dto.VoteRequest) (*dto.VoteResponse, error)
}