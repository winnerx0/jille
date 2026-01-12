package option

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Option struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name    string    `gorm:"not null"`
	PollID  uuid.UUID `gorm:"not null"`
	VoterID uuid.UUID `gorm:"not null"`
}
