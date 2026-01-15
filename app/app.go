package app

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"github.com/gofiber/swagger/v2"
	"github.com/winnerx0/jille/api/middleware"
	"github.com/winnerx0/jille/config"
	"github.com/winnerx0/jille/infra/database"
	"github.com/winnerx0/jille/infra/persistence"
	"github.com/winnerx0/jille/internal/application"
	"github.com/winnerx0/jille/internal/delivery/web"
	"github.com/winnerx0/jille/internal/domain"
	"github.com/winnerx0/jille/internal/utils"
)

type App struct {
	Config *config.Config

	Router *fiber.App

	DB *gorm.DB
}

func New(cfg *config.Config) (*App, error) {

	validator := &utils.XValidator{
		Validator: utils.Validate,
	}

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

	if err == nil {
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	}

	db.AutoMigrate(&domain.User{}, &domain.RefreshToken{}, &domain.Poll{}, &domain.Option{}, &domain.Vote{})

	app := &App{
		Config: cfg,
		Router: r,
		DB:     db,
	}

	userRepo := persistence.NewUserReposiory(db)

	pollRepo := persistence.NewPollRepository(db)

	optionRepo := persistence.NewOptionRepository(db)

	pollService := application.NewPollService(pollRepo, optionRepo)

	userService := application.NewUserService(userRepo, pollService)

	userHandler := web.NewUserHandler(userService)

	authRepo := persistence.NewAuthRepository(db)

	jwtService := application.NewJwtService(cfg.JWT_ACCESS_TOKEN_SECRET, cfg.JWT_REFRESH_TOKEN_SECRET)

	authService := application.NewAuthService(authRepo, userService, jwtService)

	authHandler := web.NewAuthHandler(authService, *validator)

	pollHandler := web.NewPollHandler(pollService, *validator)

	voteRepo := persistence.NewVoteRepository(db)

	voteservice := application.NewVoteService(voteRepo, pollRepo, optionRepo)

	voteHandler := web.NewVoteHandler(voteservice)

	app.Router.Get("/swagger/*", swagger.HandlerDefault)

	apiRouter := app.Router.Group("/api/v1")

	// auth routers

	// @Summary Register user to Jille
	// @Accept json
	// @Produce json
	// @Success 200 {object}
	// @Failure 400
	// @Router /auth/register [POST]
	apiRouter.Post("/auth/register", authHandler.RegisterUser)

	apiRouter.Post("/auth/login", authHandler.LoginUser)

	apiRouter.Post("/auth/refresh", authHandler.RefreshToken)

	// user routers

	userRouter := apiRouter.Group("/user", func(c fiber.Ctx) error {
		return middleware.JWTMiddleware(c, jwtService)
	})

	userRouter.Get("/:userID", userHandler.GetUser)

	// poll routers

	pollRouter := apiRouter.Group("/poll", func(c fiber.Ctx) error {
		return middleware.JWTMiddleware(c, jwtService)
	})

	pollRouter.Post("/create", pollHandler.CreatePoll)

	pollRouter.Post("/:pollID", pollHandler.DeletePoll)

	// vote routers

	voteRouter := apiRouter.Group("/vote", func(c fiber.Ctx) error {
		return middleware.JWTMiddleware(c, jwtService)
	})

	voteRouter.Post("/:pollId/:optionId", voteHandler.VotePoll)

	return app, err
}

func (a *App) Run() {
	fmt.Println("Server running on port", a.Config.Port)
	if err := a.Router.Listen(":" + a.Config.Port); err != nil {
		fmt.Println("Error occurred when running server", err)
	}
}
