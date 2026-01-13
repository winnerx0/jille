package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/domain"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repository.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (repo authRepository) RevokeAllTokens(ctx context.Context, userID uuid.UUID) error {

	tx := repo.db.Model(&domain.RefreshToken{}).
		Where("user_id IN (?)", userID).
		UpdateColumn("revoked", true)

	return tx.Error
}

func (repo authRepository) SaveToken(ctx context.Context, token *domain.RefreshToken) error {

	err := gorm.G[domain.RefreshToken](repo.db).Create(ctx, token)

	return err
}

func (repo authRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {

	refreshToken, err := gorm.G[domain.RefreshToken](repo.db).Where("token = ?", token).First(ctx)

	return &refreshToken, err

}