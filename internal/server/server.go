package server

import (
	"log"
	"net/http"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/cookie"
	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	handlers   *handlers
	app        *app.App
	flagConfig *config.FlagConfig
	log        *logger.Logger
}

func NewServer(app *app.App, flagConfig *config.FlagConfig, log *logger.Logger) *Server {
	handlers := newHandlers(app, flagConfig)
	return &Server{handlers: handlers, app: app, flagConfig: flagConfig, log: log}
}

func (server *Server) newRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(server.log.WithLogging())
	router.Use(middleware.CompressorMiddleware())
	router.Get("/ping", server.handlers.pingPostgresqlHandler)
	router.Post("/", server.handlers.shortenerHandler)
	router.Post("/api/shorten", server.handlers.shortenerHandlerJSON)
	router.Post("/api/shorten/batch", server.handlers.shortenerBatchHandler)
	router.Get("/{id}", server.handlers.originalHandler)
	router.Route("/api/user/urls", func(r chi.Router) {
		r.Use(cookie.CookieMiddleware())
		r.Get("/", server.handlers.urlsByIDHandler)
	})
	return router
}

func Run(server *Server) error {
	log.Println("Running server on", server.flagConfig.FlagRunAddr)
	return http.ListenAndServe(server.flagConfig.FlagRunAddr, server.newRouter())
}
