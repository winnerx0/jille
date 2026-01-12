package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/winnerx0/jille/internal/domain/option"
	"github.com/winnerx0/jille/internal/domain/poll"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string    `gorm:"unique"`
	Email    string    `gorm:"not null;unique"`
	Password string    `gorm:"not null"`
	JoinedAt time.Time `gorm:"not null"`

	Options []option.Option `gorm:"foreignKey:VoterID;references:ID"`
	Polls   []poll.Poll     `gorm:"foreignKey:UserID;references:ID"`
}
