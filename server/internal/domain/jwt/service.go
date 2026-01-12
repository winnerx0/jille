package jwt

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type jwtservice struct {
	accessTokenSecret string

	refreshTokenSecret string
}

type JWTClaims struct {
	RegisteredClaims jwt.RegisteredClaims

	// Role string `json:"role;omitempty"`

}

func NewJwtService(accessTokenSecret string, refreshTokenSecret string) Service {
	return &jwtservice{
		accessTokenSecret:  accessTokenSecret,
		refreshTokenSecret: refreshTokenSecret,
	}
}

func (j *jwtservice) GenerateAccessToken(userId string) (string, error) {

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	token, err := jwt.SignedString([]byte(j.GetAccessTokenSecretKey()))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *jwtservice) GenerateRefreshToken(userId string) (string, error) {

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	token, err := jwt.SignedString([]byte(j.GetRefreshTokenSecretKey()))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *jwtservice) VerifyAccessToken(token string) (bool, error) {

	jwtToken, err := jwt.Parse(token, func(ts *jwt.Token) (interface{}, error) {
		return []byte(j.GetAccessTokenSecretKey()), nil
	})

	fmt.Println(jwtToken)

	if err != nil {
		return false, err
	}

	return jwtToken.Valid, nil
}

func (j *jwtservice) GetTokenSubject(token string) (string, error) {

	jwtToken, err := jwt.Parse(token, func(ts *jwt.Token) (interface{}, error) {
		return []byte(j.GetAccessTokenSecretKey()), nil
	})

	if err != nil {
		return "", err
	}

	sub, err := jwtToken.Claims.GetSubject()

	if err != nil {
		return "", err
	}

	return sub, nil
}

func (j *jwtservice) VerifyRefreshToken(token string) (bool, error) {

	jwtToken, err := jwt.Parse(token, func(ts *jwt.Token) (interface{}, error) {
		return j.GetRefreshTokenSecretKey(), nil
	})

	if err != nil {
		return false, err
	}

	return jwtToken.Valid, nil
}

func (j *jwtservice) GetAccessTokenSecretKey() string {
	return j.accessTokenSecret
}

func (j *jwtservice) GetRefreshTokenSecretKey() string {
	return j.refreshTokenSecret
}
