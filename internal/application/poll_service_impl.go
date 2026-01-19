package application

import (
	"context"
	"errors"

	"fmt"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/domain"
	"github.com/winnerx0/jille/internal/utils"
)

type pollservice struct {
	repo       repository.PollRepository
	optionrepo repository.OptionRepository
	voterepo   repository.VoteRepository
}

func NewPollService(repo repository.PollRepository, optionrepo repository.OptionRepository, voterepo repository.VoteRepository) PollService {
	return &pollservice{
		repo:       repo,
		optionrepo: optionrepo,
		voterepo:   voterepo,
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

	if ctx.Value("userID").(string) != poll.UserID.String() {
		return &dto.PollViewResponse{}, utils.PollAccessDeniedError
	}

	options, err := s.optionrepo.FindOptionsByPollID(ctx, pollID)

	if err != nil {
		return &dto.PollViewResponse{}, err
	}

	var opts []dto.Option

	for _, o := range *options {

		votes := []dto.Vote{}

		for _, v := range o.Votes {
			vote := dto.Vote{
				ID:       v.ID.String(),
				UserID:   v.UserID.String(),
				PollID:   v.PollID.String(),
				OptionID: v.OptionID.String(),
			}

			votes = append(votes, vote)
		}
		option := dto.Option{
			ID:    o.ID.String(),
			Votes: votes,
			Name:  o.Name,
		}

		opts = append(opts, option)
	}

	return &dto.PollViewResponse{
		ID:        pollID.String(),
		Title:     poll.Title,
		Options:   opts,
		CreatedAt: poll.CreatedAt,
		ExpiresAt: poll.ExpiresAt,
		CreatorID: poll.UserID.String(),
	}, nil
}

func (s *pollservice) GetPoll(ctx context.Context, pollID uuid.UUID) (*dto.PollViewResponse, error) {

	poll, err := s.repo.FindPollByID(ctx, pollID)

	if err != nil {
		return &dto.PollViewResponse{}, err
	}

	options, err := s.optionrepo.FindOptionsByPollID(ctx, pollID)

	if err != nil {
		return &dto.PollViewResponse{}, err
	}

	var opts []dto.Option

	for _, o := range *options {
		option := dto.Option{
			ID:    o.ID.String(),
			Votes: []dto.Vote{}, // Hidden for public voting page
			Name:  o.Name,
		}

		opts = append(opts, option)
	}
	
	userID := ctx.Value("userID").(string)

	voted, err := s.voterepo.ExistsByPollIDAndAndUserID(ctx, pollID, uuid.MustParse(userID))

	if err != nil {
		return &dto.PollViewResponse{}, err
	}

	return &dto.PollViewResponse{
		ID:        pollID.String(),
		Title:     poll.Title,
		Options:   opts,
		CreatedAt: poll.CreatedAt,
		ExpiresAt: poll.ExpiresAt,
		CreatorID: poll.UserID.String(),
		Voted:     voted,
	}, nil
}

func (s *pollservice) GetAllPolls(ctx context.Context) (dto.ApiResponse[[]dto.PollViewResponse], error) {

	polls, err := s.repo.FindAllPolls(ctx)

	fmt.Println("polls", polls)

	if err != nil {
		return dto.ApiResponse[[]dto.PollViewResponse]{Message: "Polls reteieved successfully", Data: []dto.PollViewResponse{}}, err
	}

	var pollResponse []dto.PollViewResponse

	for _, poll := range polls {

		var opts []dto.Option

		for _, o := range poll.Options {

			votes := []dto.Vote{}

			for _, v := range o.Votes {
				vote := dto.Vote{
					ID:       v.ID.String(),
					UserID:   v.UserID.String(),
					PollID:   v.PollID.String(),
					OptionID: v.OptionID.String(),
				}

				votes = append(votes, vote)
			}

			option := dto.Option{
				ID:    o.ID.String(),
				Votes: votes,
				Name:  o.Name,
			}

			opts = append(opts, option)
		}

		response := dto.PollViewResponse{
			ID:        poll.ID.String(),
			Title:     poll.Title,
			Options:   opts,
			CreatedAt: poll.CreatedAt,
			ExpiresAt: poll.ExpiresAt,
			CreatorID: poll.UserID.String(),
		}

		pollResponse = append(pollResponse, response)
	}

	fmt.Println("poll response", pollResponse)

	return dto.ApiResponse[[]dto.PollViewResponse]{Message: "Polls reteieved successfully", Data: pollResponse}, nil
}
