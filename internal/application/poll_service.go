package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/common/dto"
)

type PollService interface {
	GetPollCount(ctx context.Context, userID uuid.UUID) (int, error)

	CreatePoll(ctx context.Context, poll *dto.CreatePollRequest) error

	DeletePoll(ctx context.Context, pollID uuid.UUID) error

	GetPollView(ctx context.Context, pollID uuid.UUID) (*dto.PollViewResponse, error)
	GetPoll(ctx context.Context, pollID uuid.UUID) (*dto.PollViewResponse, error)

	GetAllPolls(ctx context.Context) (dto.ApiResponse[[]dto.PollViewResponse], error)
}
