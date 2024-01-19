package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func LinkRouter(handlers *Handlers) chi.Router {
	router := chi.NewRouter()
	router.Post("/", handlers.ShortenerHandler)
	router.Get("/{id}", handlers.OriginalHandler)
	return router
}

func StartLinkRouter(handlers *Handlers) error {
	log.Println("Running server on", handlers.flagConfig.FlagRunAddr)
	return http.ListenAndServe(handlers.flagConfig.FlagRunAddr, LinkRouter(handlers))
}
