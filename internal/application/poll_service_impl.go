package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/domain"
	"github.com/winnerx0/jille/internal/utils"
)

type pollservice struct {
	repo       repository.PollRepository
	optionrepo repository.OptionRepository
}

func NewPollService(repo repository.PollRepository, optionrepo repository.OptionRepository) PollService {
	return &pollservice{
		repo:       repo,
		optionrepo: optionrepo,
	}
}

func (s *pollservice) GetPollCount(ctx context.Context, userID uuid.UUID) (int, error) {

	count, err := s.repo.FindUserPollCount(ctx, userID)

	if err != nil {
		return count, err
	}

	return count, nil
}

func (s *pollservice) CreatePoll(ctx context.Context, pollRequest *dto.CreatePollRequest) error {

	userID := ctx.Value("userID").(string)

	poll := &domain.Poll{
		Title:     pollRequest.Title,
		UserID:    uuid.MustParse(userID),
		ExpiresAt: pollRequest.ExpiresAt,
	}

	err := s.repo.Save(ctx, poll)

	if err != nil {
		return err
	}

	options := make([]domain.Option, len(pollRequest.Options))

	for i, option := range pollRequest.Options {
		options[i] = domain.Option{
			Name:   option,
			PollID: poll.ID,
		}
	}

	err = s.optionrepo.Save(ctx, &options)

	if err != nil {
		return err
	}

	return nil
}

func (s *pollservice) DeletePoll(ctx context.Context, pollID uuid.UUID) error {

	poll, err := s.repo.FindPollByID(ctx, pollID)

	if err != nil {
		return err
	}

	if poll.UserID != uuid.MustParse(ctx.Value("userID").(string)) {
		return errors.New("unauthorized")
	}

	err = s.repo.Delete(ctx, pollID)

	if err != nil {
		return err
	}

	return nil
}

func (s *pollservice) GetPollView(ctx context.Context, pollID uuid.UUID) (*dto.PollViewResponse, error) {

	poll, err := s.repo.FindPollByID(ctx, pollID)

	if err != nil {
		return &dto.PollViewResponse{}, err
	}

	if ctx.Value("userid").(string) != poll.UserID.String() {
		return &dto.PollViewResponse{}, utils.PollAccessDeniedError
	}

	options, err := s.optionrepo.FindOptionsByPollID(ctx, pollID)

	if err != nil {
		return &dto.PollViewResponse{}, err
	}

	var opts []dto.Option

	for _, o := range *options {
		option := dto.Option{
			ID:    o.ID.String(),
			Count: len(o.Votes),
		}

		opts = append(opts, option)
	}

	return &dto.PollViewResponse{
		ID:      pollID.String(),
		Title:   poll.Title,
		Options: opts,
	}, nil
}
