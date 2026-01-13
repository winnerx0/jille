package application

type JwtService interface {
	GenerateAccessToken(userId string) (string, error)
	GenerateRefreshToken(userId string) (string, error)
	VerifyAccessToken(token string) (bool, error)
	GetTokenSubject(token string) (string, error)
	VerifyRefreshToken(token string) (bool, error)
	GetAccessTokenSecretKey() string
	GetRefreshTokenSecretKey() string
}
