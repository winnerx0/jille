package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/winnerx0/jille/internal/domain"
)

// UserRepository Mock
type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) FindById(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepository) Save(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// PollRepository Mock
type PollRepository struct {
	mock.Mock
}

func (m *PollRepository) FindUserPollCount(ctx context.Context, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *PollRepository) Save(ctx context.Context, poll *domain.Poll) error {
	args := m.Called(ctx, poll)
	return args.Error(0)
}

func (m *PollRepository) FindPollByID(ctx context.Context, pollID uuid.UUID) (*domain.Poll, error) {
	args := m.Called(ctx, pollID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Poll), args.Error(1)
}

func (m *PollRepository) Delete(ctx context.Context, pollID uuid.UUID) error {
	args := m.Called(ctx, pollID)
	return args.Error(0)
}

// OptionRepository Mock
type OptionRepository struct {
	mock.Mock
}

func (m *OptionRepository) Save(ctx context.Context, options *[]domain.Option) error {
	args := m.Called(ctx, options)
	return args.Error(0)
}

func (m *OptionRepository) FindOptionsByPollID(ctx context.Context, pollID uuid.UUID) (*[]domain.Option, error) {


	args := m.Called(ctx, pollID)

	return args.Get(0).(*[]domain.Option), args.Error(1)
}

// AuthRepository Mock
type AuthRepository struct {
	mock.Mock
}

func (m *AuthRepository) RevokeAllTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *AuthRepository) SaveToken(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *AuthRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *AuthRepository) Delete(ctx context.Context, pollID uuid.UUID) error {
	args := m.Called(ctx, pollID)
	return args.Error(0)
}
