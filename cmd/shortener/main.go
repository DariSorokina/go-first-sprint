package main

import (
	"log"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/server"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
)

func main() {
	flagConfig := config.ParseFlags()

	var l *logger.Logger
	var err error
	if l, err = logger.CreateLogger(flagConfig.FlagLogLevel); err != nil {
		log.Fatal("Failed to create logger:", err)
	}

	storage, err := storage.SetStorage(flagConfig, l)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	app := app.NewApp(storage, l)
	serv := server.NewServer(app, flagConfig, l)

	if err := server.Run(serv); err != nil {
		panic(err)
	}
}
