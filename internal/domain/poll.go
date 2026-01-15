package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Poll struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;"`
	Title     string         `gorm:"required;not null"`
	Options   []Option       `gorm:"required;foreignKey:PollID;references:ID"`
	CreatedAt time.Time      `gorm:"required;not null"`
	UpdatedAt time.Time      `gorm:"required;not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uuid.UUID      `gorm:"required;type:uuid;not null"`
}

func (p *Poll) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
