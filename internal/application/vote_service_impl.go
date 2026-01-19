package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/utils"
)

type voteservice struct {
	repo       repository.VoteRepository
	pollrepo   repository.PollRepository
	optionrepo repository.OptionRepository
}

func NewVoteService(repo repository.VoteRepository, pollrepo repository.PollRepository, optionrepo repository.OptionRepository) VoteService {
	return &voteservice{
		repo:       repo,
		pollrepo:   pollrepo,
		optionrepo: optionrepo,
	}
}

func (s *voteservice) VotePoll(ctx context.Context, voteRequest dto.VoteRequest) (*dto.VoteResponse, error) {

	userID := ctx.Value("userID").(string)

	poll, err := s.pollrepo.FindPollByID(ctx, uuid.MustParse(voteRequest.PollID))

	if err != nil {
		return &dto.VoteResponse{}, err
	}

	if poll.ExpiresAt.Before(time.Now()) {
		return &dto.VoteResponse{}, utils.PollExpiredError
	}

	optionExists := false

	for _, option := range poll.Options {

		if option.ID == uuid.MustParse(voteRequest.OptionID) {
			optionExists = true
		}
	}

	if !optionExists {
		return &dto.VoteResponse{}, utils.OptionNotFound
	}

	err = s.repo.Vote(ctx, uuid.MustParse(voteRequest.PollID), uuid.MustParse(voteRequest.OptionID), uuid.MustParse(userID))

	if err != nil {
		return &dto.VoteResponse{}, err
	}

	return &dto.VoteResponse{
		Message: "Voted successfully",
	}, nil

}
