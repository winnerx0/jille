package auth

import (
	"context"

	"github.com/winnerx0/jille/internal/domain/user"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) Repository {
	return &authRepository{
		db: db,
	}
}

func (repo authRepository) RevokeAllTokens(ctx context.Context, email string) error {

	sub := repo.db.Model(&user.User{}).Select("id").Where("email = ?", email)
	tx := repo.db.Model(&RefreshToken{}).
		Where("user_id IN (?)", sub).
		UpdateColumn("revoked", true)

	return tx.Error
}

func (repo authRepository) SaveToken(ctx context.Context, token *RefreshToken) error {

	err := gorm.G[RefreshToken](repo.db).Create(ctx, token)

	return err
}
