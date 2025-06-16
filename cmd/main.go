package main

import (
	"github.com/paxaf/HezzlTest/config"
	"github.com/paxaf/HezzlTest/internal/app"
	"github.com/paxaf/HezzlTest/internal/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err, "Failed to load config")
	}
	app, err := app.New(cfg)
	if err != nil {
		logger.Fatal(err, "Error creating app")
	}

	if err = app.Run(); err != nil {
		logger.Fatal(err, "Error running app")
	}

	if err := app.Close(); err != nil {
		logger.Fatal(err, "Error while shutdown app")
	}
}
