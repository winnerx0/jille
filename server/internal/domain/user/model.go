package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {

	gorm.Model
	ID uuid.UUID
	Email string
	Password string
	JoinedAt time.Time
}