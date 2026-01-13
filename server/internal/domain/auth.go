package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;"`
	Token     string         `gorm:"not null"`
	ExpiresAt time.Time      `gorm:"not null"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null"`
	Revoked   bool           `gorm:"not null;default:false"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
