package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/winnerx0/jille/internal/common/dto"
	"github.com/winnerx0/jille/internal/domain/jwt"
	"github.com/winnerx0/jille/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type authservice struct {
	authrepo Repository

	userservice user.Service

	jwtservice jwt.Service
}

func NewAuthService(authrepo Repository, userservice user.Service, jwtservice jwt.Service) Service {

	return &authservice{
		authrepo:    authrepo,
		userservice: userservice,
		jwtservice:  jwtservice,
	}
}

func (s *authservice) Register(ctx context.Context, registerRequest dto.CreateUserRequest) (*dto.AuthResponse, error) {

	existingUser, err := s.userservice.GetUserByEmail(ctx, registerRequest.Email)

	if err != nil {
		return nil, err
	}

	if existingUser.Email != "" {
		return nil, errors.New("User with email already exists")
	}

	fmt.Println("user", existingUser)

	var user user.User

	user.Email = registerRequest.Email
	user.Username = registerRequest.Username

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), 10)

	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	s.userservice.CreateUser(ctx, user)

	accessToken, err := s.jwtservice.GenerateAccessToken(existingUser.ID.String())

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtservice.GenerateAccessToken(registerRequest.Email)

	if err != nil {
		return nil, err
	}

	token := RefreshToken{
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		UserID:    existingUser.ID,
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
		return &dto.AuthResponse{}, err
	}

	s.authrepo.RevokeAllTokens(ctx, loginRequest.Email)

	fmt.Println("user", existingUser)

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginRequest.Password)); err != nil {
		return nil, err
	}

	accessToken, err := s.jwtservice.GenerateAccessToken(existingUser.ID.String())

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtservice.GenerateAccessToken(loginRequest.Email)

	if err != nil {
		return nil, err
	}

	token := RefreshToken{
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
