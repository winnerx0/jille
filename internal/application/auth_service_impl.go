package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/domain"
	"github.com/winnerx0/jille/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authservice struct {
	authrepo repository.AuthRepository

	userservice UserService

	jwtservice JwtService
}

func NewAuthService(authrepo repository.AuthRepository, userservice UserService, jwtservice JwtService) AuthService {

	return &authservice{
		authrepo:    authrepo,
		userservice: userservice,
		jwtservice:  jwtservice,
	}
}

func (s *authservice) Register(ctx context.Context, registerRequest dto.CreateUserRequest) (*dto.AuthResponse, error) {

	existingUser, err := s.userservice.GetUserByEmail(ctx, registerRequest.Email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {

		} else {
			return nil, err
		}
	}

	fmt.Println("user", existingUser)
	if existingUser != nil {
		return nil, utils.UserExistsError
	}


	var user domain.User

	user.Email = registerRequest.Email
	user.Username = registerRequest.Username

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), 10)

	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	s.userservice.CreateUser(ctx, &user)

	fmt.Println("user id", user.ID)

	accessToken, err := s.jwtservice.GenerateAccessToken(user.ID.String())

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtservice.GenerateAccessToken(registerRequest.Email)

	if err != nil {
		return nil, err
	}

	token := domain.RefreshToken{
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		UserID:    user.ID,
	}

	s.authrepo.SaveToken(ctx, &token)

	return &dto.AuthResponse{
		Message: "Regtration successful",
		AuthTokens: dto.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *authservice) Login(ctx context.Context, loginRequest dto.LoginUserRequest) (*dto.AuthResponse, error) {

	existingUser, err := s.userservice.GetUserByEmail(ctx, loginRequest.Email)

	if err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, utils.UserNotFoundError
		} else {
			return nil, err
		}
	}

	s.authrepo.RevokeAllTokens(ctx, existingUser.ID)

	fmt.Println("user", existingUser)

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginRequest.Password)); err != nil {
		return nil, errors.New("Invalid password")
	}

	accessToken, err := s.jwtservice.GenerateAccessToken(existingUser.ID.String())

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtservice.GenerateAccessToken(loginRequest.Email)

	if err != nil {
		return nil, err
	}

	token := domain.RefreshToken{
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		UserID:    existingUser.ID,
	}

	s.authrepo.SaveToken(ctx, &token)

	return &dto.AuthResponse{
		Message: "Login successful",
		AuthTokens: dto.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *authservice) RefreshToken(ctx context.Context, refreshTokenRequest dto.RefreshTokenRequest) (*dto.AuthResponse, error) {

	existingToken, err := s.authrepo.FindByToken(ctx, refreshTokenRequest.RefreshToken)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.TokenNotFoundError
		} else {
			return nil, err
		}
	}

	if existingToken.ExpiresAt.Before(time.Now()) {
		return nil, utils.TokenExpiredError
	}

	accessToken, err := s.jwtservice.GenerateAccessToken(existingToken.UserID.String())

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtservice.GenerateAccessToken(existingToken.UserID.String())

	if err != nil {
		return nil, err
	}

	err = s.authrepo.RevokeAllTokens(ctx, existingToken.UserID)

	if err != nil {
		return nil, err
	}

	token := domain.RefreshToken{
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		UserID:    existingToken.UserID,
	}

	s.authrepo.SaveToken(ctx, &token)

	return &dto.AuthResponse{
		Message: "Refresh token successful",
		AuthTokens: dto.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}