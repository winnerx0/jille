package main

import (
	"log"
	
	"github.com/joho/godotenv"
	"github.com/winnerx0/jille/app"
	"github.com/winnerx0/jille/config"
)

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
