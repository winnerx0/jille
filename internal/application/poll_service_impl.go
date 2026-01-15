package application

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/domain"
)

type pollservice struct {
	repo        repository.PollRepository
	optionrepo  repository.OptionRepository
	userservice UserService
}

func NewPollService(repo repository.PollRepository, optionrepo repository.OptionRepository) PollService {
	return &pollservice{
    repo: repo,
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
		Title: pollRequest.Title,
		UserID: uuid.MustParse(userID),
	}

	err := s.repo.Save(ctx, poll)

	if err != nil {
		return err
	}


	options := make([]domain.Option, len(pollRequest.Options))

	for i, option := range pollRequest.Options {
		options[i] = domain.Option{
			Name: option,
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