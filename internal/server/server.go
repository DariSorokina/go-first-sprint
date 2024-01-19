package server

import (
	"log"
	"net/http"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	handlers   *handlers
	app        *app.App
	flagConfig *config.FlagConfig
}

func NewServer(app *app.App, flagConfig *config.FlagConfig) *Server {
	handlers := newHandlers(app, flagConfig)
	return &Server{handlers: handlers, app: app, flagConfig: flagConfig}
}

func (server *Server) newRouter(handlers *handlers) chi.Router {
	router := chi.NewRouter()
	router.Post("/", handlers.shortenerHandler)
	router.Get("/{id}", handlers.originalHandler)
	return router
}

func Run(server *Server) error {
	log.Println("Running server on", server.flagConfig.FlagRunAddr)
	return http.ListenAndServe(server.flagConfig.FlagRunAddr, server.newRouter(server.handlers))
}
