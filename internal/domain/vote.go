package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vote struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_user_poll"`
	PollID    uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_user_poll"`
	OptionID  uuid.UUID      `gorm:"type:uuid;not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (v *Vote) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return
}
