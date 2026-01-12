package auth

import (
	"context"
)

type Repository interface {

	RevokeAllTokens(ctx context.Context, email string) error

	SaveToken(ctx context.Context, token *RefreshToken) error
}