package server

import (
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

func NewServer(app *app.App, flagConfig *config.FlagConfig, l *logger.Logger) *Server {
	handlers := newHandlers(app, flagConfig, l)
	return &Server{handlers: handlers, app: app, flagConfig: flagConfig, log: l}
}

func (server *Server) newRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(server.log.WithLogging())
	router.Use(middleware.CompressorMiddleware())
	router.Get("/ping", server.handlers.pingPostgresqlHandler)
	router.Get("/{id}", server.handlers.originalHandler)
	router.Route("/", func(r chi.Router) {
		r.Use(cookie.CookieMiddleware())
		r.Post("/", server.handlers.shortenerHandler)
		r.Post("/api/shorten", server.handlers.shortenerHandlerJSON)
		r.Post("/api/shorten/batch", server.handlers.shortenerBatchHandler)
		r.Get("/api/user/urls", server.handlers.urlsByIDHandler)
		r.Delete("/api/user/urls", server.handlers.deleteURLsHandler)
	})
	return router
}

func Run(server *Server) error {
	server.log.Sugar().Infof("Running server on %s", server.flagConfig.FlagRunAddr)
	return http.ListenAndServe(server.flagConfig.FlagRunAddr, server.newRouter())
}
