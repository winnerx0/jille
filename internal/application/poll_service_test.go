package application

import (
	"context"
	"testing"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/winnerx0/jille/internal/application/repository/mocks"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/domain"
)

func TestGetPollCount(t *testing.T) {
	mockRepo := new(mocks.PollRepository)
	mockOptionRepo := new(mocks.OptionRepository)
	service := NewPollService(mockRepo, mockOptionRepo)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("FindUserPollCount", ctx, userID).Return(3, nil)

	count, err := service.GetPollCount(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	mockRepo.AssertExpectations(t)
}

func TestCreatePoll_Success(t *testing.T) {
	mockRepo := new(mocks.PollRepository)
	mockOptionRepo := new(mocks.OptionRepository)
	service := NewPollService(mockRepo, mockOptionRepo)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), "userID", userID.String())

	pollRequest := &dto.CreatePollRequest{
		Title:   "Test Poll",
		Options: []string{"Option 1", "Option 2"},
	}

	mockRepo.On("Save", ctx, mock.MatchedBy(func(p *domain.Poll) bool {
		return p.Title == "Test Poll" && p.UserID == userID
	})).Return(nil)

	mockOptionRepo.On("Save", ctx, mock.MatchedBy(func(opts *[]domain.Option) bool {
		return len(*opts) == 2
	})).Return(nil)

	err := service.CreatePoll(ctx, pollRequest)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockOptionRepo.AssertExpectations(t)
}

func TestDeletePoll(t *testing.T) {
	mockRepo := new(mocks.PollRepository)
	mockOptionRepo := new(mocks.OptionRepository)
	service := NewPollService(mockRepo, mockOptionRepo)

	ctx := context.Background()
	pollID := uuid.New()

	mockRepo.On("Delete", ctx, pollID).Return(nil)

	err := service.DeletePoll(ctx, pollID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeletePoll_Fail(t *testing.T) {
	mockRepo := new(mocks.PollRepository)
	mockOptionRepo := new(mocks.OptionRepository)
	service := NewPollService(mockRepo, mockOptionRepo)

	ctx := context.Background()
	pollID := uuid.New()

	mockRepo.On("Delete", ctx, pollID).Return(errors.New("Poll not found"))

	err := service.DeletePoll(ctx, pollID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}