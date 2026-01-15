package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/application/repository"
	"github.com/winnerx0/jille/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserReposiory(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) FindById(ctx context.Context, userId uuid.UUID) (domain.User, error) {

	user, err := gorm.G[domain.User](repo.db).Where("id = ?", userId).First(ctx)

	return user, err
}

func (repo *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {

	exists, err := gorm.G[bool](repo.db).Raw("SELECT COUNT(u) > 0 FROM users u WHERE u.email = ?", email).First(ctx)

	return exists, err
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {

	user, err := gorm.G[domain.User](repo.db).Where("email = ?", email).First(ctx)

	return user, err
}

func (repo *userRepository) Save(ctx context.Context, user *domain.User) error {

	err := gorm.G[domain.User](repo.db).Create(ctx, user)

	return err
}
