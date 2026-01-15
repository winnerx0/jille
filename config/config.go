package config

import (
	"errors"
	"os"

	"github.com/winnerx0/jille/infra/database"
)

type Config struct {
	Port                     string
	JWT_ACCESS_TOKEN_SECRET  string
	JWT_REFRESH_TOKEN_SECRET string
	DBConfig                 database.DBConfig
}

func Load() (*Config, error) {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	jwt_access_token_secret := os.Getenv("JWT_ACCESS_TOKEN_SECRET")

	if jwt_access_token_secret == "" {
		return nil, errors.New("JWT Access Token Secret Required")
	}

	jwt_refresh_token_secret := os.Getenv("JWT_REFRESH_TOKEN_SECRET")

	if jwt_refresh_token_secret == "" {
		return nil, errors.New("JWT Refresh Token Secret Required")
	}

	db_host := os.Getenv("DB_HOST")
	if db_host == "" {
		return nil, errors.New("DB Host Required")
	}

	db_port := os.Getenv("DB_PORT")
	if db_port == "" {
		return nil, errors.New("DB Port Required")
	}

	db_user := os.Getenv("DB_USER")
	if db_user == "" {
		return nil, errors.New("DB User Required")
	}

	db_password := os.Getenv("DB_PASSWORD")
	if db_password == "" {
		return nil, errors.New("DB Password Required")
	}

	db_name := os.Getenv("DB_NAME")
	if db_name == "" {
		return nil, errors.New("DB Name Required")
	}

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	timeZone := os.Getenv("DB_TIMEZONE")
	if timeZone == "" {
		timeZone = "UTC"
	}

	cfg := &Config{
		Port:                     port,
		JWT_ACCESS_TOKEN_SECRET:  jwt_access_token_secret,
		JWT_REFRESH_TOKEN_SECRET: jwt_refresh_token_secret,
		DBConfig: database.DBConfig{
			Host:     db_host,
			Port:     db_port,
			User:     db_user,
			Password: db_password,
			Name:     db_name,
			SSLMode:  sslMode,
			TimeZone: timeZone,
		},
	}

	return cfg, nil
}
