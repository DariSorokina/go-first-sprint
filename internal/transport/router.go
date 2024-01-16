package transport

import (
	"log"
	"net/http"

	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/go-chi/chi/v5"
)

func LinkRouter() chi.Router {
	router := chi.NewRouter()
	router.Post("/", ShortenerHandler)
	router.Get("/{id}", OriginalHandler)
	return router
}

func RunLinkRouter() error {
	log.Println("Running server on", config.FlagRunAddr)
	return http.ListenAndServe(config.FlagRunAddr, LinkRouter())
}
