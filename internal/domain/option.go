package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Option struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;"`
	Name      string         `gorm:"not null"`
	PollID    uuid.UUID      `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Votes []Vote `gorm:"foreignKey:OptionID;references:ID"`
}

func (o *Option) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
