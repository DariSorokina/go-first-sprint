package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/models"
	"github.com/go-chi/chi/v5"
)

type originalURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type shortURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type handlers struct {
	app        *app.App
	flagConfig *config.FlagConfig
}

func newHandlers(app *app.App, flagConfig *config.FlagConfig) *handlers {
	return &handlers{app: app, flagConfig: flagConfig}
}

func (handlers *handlers) pingPostgresqlHandler(res http.ResponseWriter, req *http.Request) {
	err := handlers.app.PingPostgresql()
	if err != nil {
		http.Error(res, "Storage connection failed", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (handlers *handlers) shortenerHandler(res http.ResponseWriter, req *http.Request) {
	var response string
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Bad request body", http.StatusBadRequest)
		return
	}

	shortenedURL := handlers.app.ToShortenURL(string(requestBody))

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

func (handlers *handlers) originalHandler(res http.ResponseWriter, req *http.Request) {
	idValue := chi.URLParam(req, "id")
	correspondingURL := handlers.app.ToOriginalURL(idValue)
	res.Header().Set("Location", correspondingURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (handlers *handlers) shortenerHandlerJSON(res http.ResponseWriter, req *http.Request) {
	var request models.Request
	var response models.Response
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	shortenedURL := handlers.app.ToShortenURL(string(request.OriginalURL))

	response.ShortenURL, err = url.JoinPath(handlers.flagConfig.FlagBaseURL, shortenedURL)
	if err != nil {
		http.Error(res, "Bad URL path provided", http.StatusInternalServerError)
		log.Println("Failed to join provided URL path with short URL", err)
		return
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}

func (handlers *handlers) shortenerBatchHandler(res http.ResponseWriter, req *http.Request) {
	var input []originalURL
	var output []shortURL
	var response string

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Bad request body", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(requestBody, &input); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	for _, inputSample := range input {
		shortenedURL := handlers.app.ToShortenURL(inputSample.OriginalURL)

		response, err = url.JoinPath(handlers.flagConfig.FlagBaseURL, shortenedURL)
		if err != nil {
			http.Error(res, "Bad URL path provided", http.StatusInternalServerError)
			log.Println("Failed to join provided URL path with short URL", err)
			return
		}

		output = append(output, shortURL{CorrelationID: inputSample.CorrelationID, ShortURL: response})
	}

	resp, err := json.Marshal(output)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}
