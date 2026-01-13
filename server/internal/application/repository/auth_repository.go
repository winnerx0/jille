package repository

import (
	"context"

	"github.com/winnerx0/jille/internal/domain"
	"github.com/google/uuid"
)

type AuthRepository interface {

	RevokeAllTokens(ctx context.Context, userID uuid.UUID) error

	// RevokeToken(ctx context.Context, token string) error

	SaveToken(ctx context.Context, token *domain.RefreshToken) error

	FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
}