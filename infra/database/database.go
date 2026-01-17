package database

import (
	"fmt"

	"github.com/winnerx0/jille/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

func (db *DBConfig) New() (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		db.Host,
		db.User,
		db.Password,
		db.Name,
		db.Port,
		db.SSLMode,
		db.TimeZone,
	)

	database, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		return nil, err
	}
	
	return database, nil
}

var Models = []any{
	&domain.User{},
	&domain.RefreshToken{},
	&domain.Poll{},
	&domain.Option{},
	&domain.Vote{},
}