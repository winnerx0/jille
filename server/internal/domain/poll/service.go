package poll

import (
	"context"
	// "errors"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/domain/user"
)

type pollservice struct {
	repo        Repository
	userservice user.Service
}

func NewPollService(repo Repository) *pollservice {
	return &pollservice{
		repo: repo,
	}
}

func (s *pollservice) GetPollCount(ctx context.Context, userID uuid.UUID) (int, error) {

	// _, err := s.userservice.GetUserById(ctx, userID)

	// if err != nil {
	// 	return 0, err
	// }

	// if user.Email == "" {
	// 	return 0, errors.New("User not found")
	// }

	count, err := s.repo.FindUserPollCount(ctx, userID)

	if err != nil {
		return count, err
	}

	return count, nil
}
