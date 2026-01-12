package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/winnerx0/jille/api/middleware"
	"github.com/winnerx0/jille/config"
	"github.com/winnerx0/jille/infra/database"
	"github.com/winnerx0/jille/internal/domain/auth"
	"github.com/winnerx0/jille/internal/domain/jwt"
	"github.com/winnerx0/jille/internal/domain/option"
	"github.com/winnerx0/jille/internal/domain/poll"
	"github.com/winnerx0/jille/internal/domain/user"
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

	db.AutoMigrate(&user.User{}, &auth.RefreshToken{}, &poll.Poll{}, &option.Option{})

	app := &App{
		Config: cfg,
		Router: r,
		DB:     db,
	}

	userRepo := user.NewUserReposiory(db)

	pollRepo := poll.NewPollRepository(db)

	pollService := poll.NewPollService(pollRepo)

	userService := user.NewUserService(userRepo, pollService)

	userHandler := user.NewUserHandler(userService)

	authRepo := auth.NewAuthRepository(db)

	jwtservice := jwt.NewJwtService(cfg.JWT_ACCESS_TOKEN_SECRET, cfg.JWT_REFRESH_TOKEN_SECRET)

	authService := auth.NewAuthService(authRepo, userService, jwtservice)

	authHandler := auth.NewAuthHandler(authService)

	// auth handlers
	app.Router.Post("/api/v1/auth/login", authHandler.LoginUser)

	// user handlers
	userRouter := app.Router.Group("/api/v1/user", func(c *fiber.Ctx) error {
		return middleware.JWTMiddleware(c, jwtservice)
	})

	userRouter.Get("/:userID", userHandler.GetUser)

	return app, err
}

func (a *App) Run() {
	fmt.Println("Server running on port", a.Config.Port)
	if err := a.Router.Listen(":" + a.Config.Port); err != nil {
		fmt.Println("Error occurred when running server", err)
	}
}
