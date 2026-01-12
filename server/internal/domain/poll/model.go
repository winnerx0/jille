package poll

import (
	"time"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/domain/option"
	"gorm.io/gorm"
)

type Poll struct {
	gorm.Model

	ID        uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title     string          `gorm:"not null"`
	Options   []option.Option `gorm:"foreignKey:PollID;references:ID"`
	CreatedAt time.Time       `gorm:"not null"`
	UserID    uuid.UUID       `gorm:"not null"`
}
