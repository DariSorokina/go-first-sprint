package main

import (
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/server"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
)

func main() {
	flagConfig := config.ParseFlags()
	urlMap := storage.NewURL()
	handlers := server.NewHandlers(urlMap, flagConfig)
	if err := server.StartLinkRouter(handlers); err != nil {
		panic(err)
	}
}
