package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string    `gorm:"unique"`
	Email    string    `gorm:"not null;unique"`
	Password string    `gorm:"not null"`
	JoinedAt time.Time `gorm:"not null"`
}
