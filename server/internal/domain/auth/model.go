package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model

	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Token     string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Revoked   bool      `gorm:"not null;default:false"`
}

 