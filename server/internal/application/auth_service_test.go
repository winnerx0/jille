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
	"gorm.io/gorm"
)

// MockUserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserById(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*dto.UserAuthView, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserAuthView), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockJwtService
type MockJwtService struct {
	mock.Mock
}

func (m *MockJwtService) GenerateAccessToken(userId string) (string, error) {
	args := m.Called(userId)
	return args.String(0), args.Error(1)
}

func (m *MockJwtService) GenerateRefreshToken(userId string) (string, error) {
	args := m.Called(userId)
	return args.String(0), args.Error(1)
}

func (m *MockJwtService) VerifyAccessToken(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}

func (m *MockJwtService) GetTokenSubject(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockJwtService) VerifyRefreshToken(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}

func (m *MockJwtService) GetAccessTokenSecretKey() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockJwtService) GetRefreshTokenSecretKey() string {
	args := m.Called()
	return args.String(0)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(mocks.AuthRepository)
	mockUserService := new(MockUserService)
	mockJwtService := new(MockJwtService)
	service := NewAuthService(mockRepo, mockUserService, mockJwtService)

	ctx := context.Background()
	req := dto.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	}

	// Mock GetUserByEmail returns record not found (User does not exist)
	mockUserService.On("GetUserByEmail", ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Mock CreateUser
	mockUserService.On("CreateUser", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*domain.User)
		user.ID = uuid.New() // Simulate ID generation
	})

	// Mock JwtService
	mockJwtService.On("GenerateAccessToken", mock.AnythingOfType("string")).Return("access_token", nil).Once()
	mockJwtService.On("GenerateRefreshToken", mock.AnythingOfType("string")).Return("refresh_token", nil).Once()

	// Mock SaveToken
	mockRepo.On("SaveToken", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	resp, err := service.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "access_token", resp.AuthTokens.AccessToken)
	assert.Equal(t, "refresh_token", resp.AuthTokens.RefreshToken) 
	mockUserService.AssertExpectations(t)
	mockJwtService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(mocks.AuthRepository)
	mockUserService := new(MockUserService)
	mockJwtService := new(MockJwtService)
	service := NewAuthService(mockRepo, mockUserService, mockJwtService)

	ctx := context.Background()
	req := dto.LoginUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Create a user with hashed password (using bcrypt cost 4 for speed in tests if possible, but here we just mock)
	// We need actual valid hash because Login calls bcrypt.CompareHashAndPassword
	// $2a$10$X... is a valid bcrypt hash for "password123"
	// But calculating it might be slow or tricky.
	// The service calls `bcrypt.CompareHashAndPassword`.
	// We can pre-calculate a hash for "password123".
	hashedPassword := "$2a$10$89.invalid.hash.but.needs.to.work.if.we.could.mock.bcrypt"
	// Actually we cannot mock bcrypt package calls. We must provide a valid hash.
	// Valid hash for "password123" (cost 10)
	// $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
	hashedPassword = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

	user := &dto.UserAuthView{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: hashedPassword,
	}

	mockUserService.On("GetUserByEmail", ctx, req.Email).Return(user, nil)
	mockRepo.On("RevokeAllTokens", ctx, user.ID).Return(nil)

	// Implementation calls GenerateAccessToken for both access and refresh
	mockJwtService.On("GenerateAccessToken", mock.AnythingOfType("string")).Return("token", nil).Twice()
	mockRepo.On("SaveToken", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	resp, err := service.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockUserService.AssertExpectations(t)
	mockJwtService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	mockRepo := new(mocks.AuthRepository)
	mockUserService := new(MockUserService)
	mockJwtService := new(MockJwtService)
	service := NewAuthService(mockRepo, mockUserService, mockJwtService)

	ctx := context.Background()
	req := dto.RefreshTokenRequest{
		RefreshToken: "valid_refresh_token",
	}

	existingToken := &domain.RefreshToken{
		ID:        uuid.New(),
		Token:     req.RefreshToken,
		UserID:    uuid.New(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	mockRepo.On("FindByToken", ctx, req.RefreshToken).Return(existingToken, nil)
	mockJwtService.On("GenerateAccessToken", existingToken.UserID.String()).Return("new_access_token", nil).Twice()
	mockRepo.On("RevokeAllTokens", ctx, existingToken.UserID).Return(nil)
	mockRepo.On("SaveToken", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	resp, err := service.RefreshToken(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "new_access_token", resp.AuthTokens.AccessToken)

	mockRepo.AssertExpectations(t)
}
