package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/winnerx0/jille/internal/common"
	"github.com/winnerx0/jille/internal/common/dto"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, loginRequest dto.LoginUserRequest) (*dto.AuthResponse, error)
}

type JwtService interface {
	GenerateAccessToken(userId string) (string, error)
	GenerateRefreshToken(userId string) (string, error)
	VerifyAccessToken(token string) (bool, error)
	GetTokenSubject(token string) (string, error)
	VerifyRefreshToken(token string) (bool, error)
	GetAccessTokenSecretKey() string
	GetRefreshTokenSecretKey() string
}

type authservice struct {
	authrepo Repository

	userservice common.UserService

	jwtservice JwtService
}

func NewAuthService(authrepo Repository, userservice common.UserService, jwtservice JwtService) *authservice {

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
