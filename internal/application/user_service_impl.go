package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/domain"
)

type userservice struct {
	userRepo    repository.UserRepository
	pollService PollService
}

func NewUserService(userRepo repository.UserRepository, pollservice PollService) UserService {
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
		ID:        user.ID,
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

func (s *userservice) CreateUser(ctx context.Context, user *domain.User) error {

	err := s.userRepo.Save(ctx, user)

	return err
}
