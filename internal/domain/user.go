package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;"`
	Username  string         `gorm:"not null"`
	Email     string         `gorm:"not null;unique"`
	Password  string         `gorm:"not null"`
	JoinedAt  time.Time      `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Polls         []Poll         `gorm:"foreignKey:UserID;references:ID"`
	Votes         []Vote         `gorm:"foreignKey:UserID;references:ID"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID;references:ID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
