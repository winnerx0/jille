package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/winnerx0/jille/app"
	"github.com/winnerx0/jille/config"
	_ "github.com/winnerx0/jille/docs"
)

// @title Jille Backend Server API
// @version 0.1
// @description This is a the documenetation for Jille API
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9000
// @BasePath /
func main() {

	godotenv.Load(".env")

	config, err := config.Load()

	if err != nil {
		log.Fatal("Failed to load config ", err.Error())
	}

	app, err := app.New(config)

	if err != nil {
		log.Fatal(err.Error())
	}

	app.Run()
}
