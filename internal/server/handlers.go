package server

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	urlMap     *storage.URL
	flagConfig *config.FlagConfig
}

func NewHandlers(urlMap *storage.URL, flagConfig *config.FlagConfig) *Handlers {
	return &Handlers{urlMap: urlMap, flagConfig: flagConfig}
}

func (handlers *Handlers) ShortenerHandler(res http.ResponseWriter, req *http.Request) {
	var response string
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Bad request body", http.StatusBadRequest)
		return
	}

	shortenedURL := app.ToSortenURL(handlers.urlMap, string(requestBody))

	response, err = url.JoinPath(handlers.flagConfig.FlagBaseURL, shortenedURL)
	if err != nil {
		http.Error(res, "Bad URL path provided", http.StatusInternalServerError)
		log.Println("Failed to join provided URL path with short URL", err)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(response))
}

func (handlers *Handlers) OriginalHandler(res http.ResponseWriter, req *http.Request) {
	idValue := chi.URLParam(req, "id")
	correspondingURL := app.ToOriginalURL(handlers.urlMap, idValue)
	res.Header().Set("Location", correspondingURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
