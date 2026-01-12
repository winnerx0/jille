package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/common/dto"
)

type pollservice interface {
	GetPollCount(ctx context.Context, userID uuid.UUID) (int, error)
}

type userservice struct {
	userRepo    Repository
	pollService pollservice
}

type UserService interface {
	GetUserById(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	GetUserByEmail(ctx context.Context, email string) (*dto.UserAuthView, error)
}

func NewUserService(userRepo Repository, pollservice pollservice) UserService {
	return &userservice{
		userRepo:    userRepo,
		pollService: pollservice,
	}
}

func (s *userservice) GetUserById(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {

	user, err := s.userRepo.FindById(ctx, userID)

	if err != nil {
		return nil, err
	}


	if user.Email == "" {
		return nil, errors.New("User not found")
	}

	pollCount, err := s.pollService.GetPollCount(ctx, userID)

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID: user.ID,
		Email:     user.Email,
		PollCount: pollCount,
	}, nil
}

func (s *userservice) ExistsByEmail(ctx context.Context, email string) (bool, error) {

	exists, err := s.userRepo.ExistsByEmail(ctx, email)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *userservice) GetUserByEmail(ctx context.Context, email string) (*dto.UserAuthView, error) {

	user, err := s.userRepo.FindByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return &dto.UserAuthView{
		Email:    user.Email,
		Password: user.Password,
		ID:       user.ID,
	}, nil
}
