package auth

import (
	"context"
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
