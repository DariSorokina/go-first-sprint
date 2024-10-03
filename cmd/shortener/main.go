package main

import (
	"fmt"
	"log"

	"github.com/DariSorokina/go-first-sprint/internal/app"
	"github.com/DariSorokina/go-first-sprint/internal/config"
	"github.com/DariSorokina/go-first-sprint/internal/logger"
	"github.com/DariSorokina/go-first-sprint/internal/server"
	"github.com/DariSorokina/go-first-sprint/internal/storage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func init() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
}

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

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
