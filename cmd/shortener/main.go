package main

import (
	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/server"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
)

func main() {
	flagConfig := config.ParseFlags()
	storage := storage.NewStorage()
	app := app.NewApp(storage)
	serv := server.NewServer(app, flagConfig)
	if err := server.Run(serv); err != nil {
		panic(err)
	}
}
