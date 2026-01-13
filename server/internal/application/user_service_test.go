package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/winnerx0/jille/internal/application/repository/mocks"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/domain"
)

func TestGetUserById_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	mockPollService := new(MockPollService)
	service := NewUserService(mockRepo, mockPollService)

	ctx := context.Background()
	userID := uuid.New()

	expectedUser := domain.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		JoinedAt: time.Now(),
	}

	mockRepo.On("FindById", ctx, userID).Return(expectedUser, nil)
	mockPollService.On("GetPollCount", ctx, userID).Return(5, nil)

	resp, err := service.GetUserById(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, resp.Email)
	assert.Equal(t, 5, resp.PollCount)
	mockRepo.AssertExpectations(t)
	mockPollService.AssertExpectations(t)
}

func TestExistsByEmail(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := NewUserService(mockRepo, nil)

	ctx := context.Background()
	email := "test@example.com"

	mockRepo.On("ExistsByEmail", ctx, email).Return(true, nil)

	exists, err := service.ExistsByEmail(ctx, email)

	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := NewUserService(mockRepo, nil)

	ctx := context.Background()
	email := "test@example.com"
	expectedUser := domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: "hashedpassword",
	}

	mockRepo.On("FindByEmail", ctx, email).Return(expectedUser, nil)

	resp, err := service.GetUserByEmail(ctx, email)

	assert.NoError(t, err)
	assert.Equal(t, email, resp.Email)
	assert.Equal(t, "hashedpassword", resp.Password)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := NewUserService(mockRepo, nil)

	ctx := context.Background()
	user := &domain.User{
		Username: "newuser",
		Email:    "new@example.com",
	}

	mockRepo.On("Save", ctx, user).Return(nil)

	err := service.CreateUser(ctx, user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// MockPollService local definition for testing
type MockPollService struct {
	mock.Mock
}

func (m *MockPollService) CreatePoll(ctx context.Context, pollRequest *dto.CreatePollRequest) error {
	args := m.Called(ctx, pollRequest)
	return args.Error(0)
}

func (m *MockPollService) GetPollCount(ctx context.Context, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockPollService) DeletePoll(ctx context.Context, pollID uuid.UUID) error {
	args := m.Called(ctx, pollID)
	return args.Error(0)
}