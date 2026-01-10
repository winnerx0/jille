package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/winnerx0/jille/config"
	"github.com/winnerx0/jille/infra/database"
	"github.com/winnerx0/jille/internal/domain/user"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config

	Router *fiber.App

	DB *gorm.DB
}

func New(cfg *config.Config) (*App, error) {

	r := fiber.New()

	database := &database.DBConfig{
		Host:     cfg.DBConfig.Host,
		Port:     cfg.DBConfig.Port,
		User:     cfg.DBConfig.User,
		Password: cfg.DBConfig.Password,
		Name:     cfg.DBConfig.Name,
		SSLMode:  cfg.DBConfig.SSLMode,
		TimeZone: cfg.DBConfig.TimeZone,
	}

	db, err := database.New()

	db.AutoMigrate(&user.User{})

	app := &App{
		Config: cfg,
		Router: r,
		DB:     db,
	}

	return app, err
}

func (a *App) Run() {

	fmt.Println("Server running on port", a.Config.Port)
	if err := a.Router.Listen(":" + a.Config.Port); err != nil {
		fmt.Println("Error occurred when running server", err)
	}
}
