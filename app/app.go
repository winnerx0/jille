package app

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/gorm"

	"github.com/winnerx0/jille/api/middleware"
	"github.com/winnerx0/jille/config"
	"github.com/winnerx0/jille/infra/database"
	"github.com/winnerx0/jille/infra/persistence"
	"github.com/winnerx0/jille/internal/application"
	"github.com/winnerx0/jille/internal/delivery/web"
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

	dbConfig := &database.DBConfig{
		Host:     cfg.DBConfig.Host,
		Port:     cfg.DBConfig.Port,
		User:     cfg.DBConfig.User,
		Password: cfg.DBConfig.Password,
		Name:     cfg.DBConfig.Name,
		SSLMode:  cfg.DBConfig.SSLMode,
		TimeZone: cfg.DBConfig.TimeZone,
	}

	db, err := dbConfig.New()
	if err != nil {
		log.Fatal("Error connecting to database", err.Error())
	}

	// if err := db.AutoMigrate(database.Models...); err != nil {
	// 		log.Fatal("Error migrating dataabse tables", err.Error())
	// }

	app := &App{
		Config: cfg,
		Router: r,
		DB:     db,
	}

	app.Router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://jille.vercel.app"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	userRepo := persistence.NewUserReposiory(db)

	pollRepo := persistence.NewPollRepository(db)

	optionRepo := persistence.NewOptionRepository(db)

	voteRepo := persistence.NewVoteRepository(db)

	pollService := application.NewPollService(pollRepo, optionRepo, voteRepo)

	userService := application.NewUserService(userRepo, pollService)

	userHandler := web.NewUserHandler(userService)

	authRepo := persistence.NewAuthRepository(db)

	jwtService := application.NewJwtService(cfg.JWT_ACCESS_TOKEN_SECRET, cfg.JWT_REFRESH_TOKEN_SECRET)

	authService := application.NewAuthService(authRepo, userService, jwtService)

	authHandler := web.NewAuthHandler(authService, *validator)

	pollHandler := web.NewPollHandler(pollService, *validator)

	broker := utils.NewBroker()
	broker.Start()

	voteservice := application.NewVoteService(voteRepo, pollRepo, optionRepo)

	voteHandler := web.NewVoteHandler(voteservice)

	apiRouter := app.Router.Group("/api/v1")

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

	pollRouter.Get("/all", pollHandler.GetAllPolls)

	pollRouter.Get("/view/:pollID", pollHandler.GetPollView)

	pollRouter.Get("/:pollID", pollHandler.GetPoll)

	// vote routers
	voteRouter := apiRouter.Group("/vote", func(c fiber.Ctx) error {
		return middleware.JWTMiddleware(c, jwtService)
	})

	voteRouter.Post("/", voteHandler.VotePoll(*broker))

	apiRouter.Get("/sse", web.SseHandler(broker))

	return app, err
}

func (a *App) Run() {
	fmt.Println("Server running on port", a.Config.Port)
	if err := a.Router.Listen(":" + a.Config.Port); err != nil {
		fmt.Println("Error occurred when running server", err)
	}
}
