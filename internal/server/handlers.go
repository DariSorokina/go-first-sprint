package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/models"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
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
	err := handlers.app.Ping()
	if err != nil {
		http.Error(res, "Storage connection failed", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (handlers *handlers) originalHandler(res http.ResponseWriter, req *http.Request) {
	idValue := chi.URLParam(req, "id")
	correspondingURL, getOriginalErr := handlers.app.ToOriginalURL(idValue)
	fmt.Println("--------------")
	fmt.Println(correspondingURL)
	if errors.Is(getOriginalErr, storage.ErrDeletedURL) {
		res.Header().Set("Location", correspondingURL)
		res.WriteHeader(http.StatusGone)
	} else {
		res.Header().Set("Location", correspondingURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}

}

func (handlers *handlers) shortenerHandler(res http.ResponseWriter, req *http.Request) {
	var response string
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Bad request body", http.StatusBadRequest)
		return
	}

	userID := req.Header.Get("ClientID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(err)
	}

	shortenedURL, errShortURL := handlers.app.ToShortenURL(string(requestBody), userIDInt)
	response, err = url.JoinPath(handlers.flagConfig.FlagBaseURL, shortenedURL)
	if err != nil {
		http.Error(res, "Bad URL path provided", http.StatusInternalServerError)
		log.Println("Failed to join provided URL path with short URL", err)
		return
	}

	res.Header().Set("content-type", "text/plain")

	if errors.Is(errShortURL, storage.ErrShortURLAlreadyExist) {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}

	res.Write([]byte(response))
}

func (handlers *handlers) shortenerHandlerJSON(res http.ResponseWriter, req *http.Request) {
	var request models.Request
	var response models.Response

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(requestBody, &request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	userID := req.Header.Get("ClientID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(err)
	}

	shortenedURL, errShortURL := handlers.app.ToShortenURL(string(request.OriginalURL), userIDInt)

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

	if errors.Is(errShortURL, storage.ErrShortURLAlreadyExist) {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}

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

	userID := req.Header.Get("ClientID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(err)
	}

	for _, inputSample := range input {
		shortenedURL, _ := handlers.app.ToShortenURL(inputSample.OriginalURL, userIDInt)

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

func (handlers *handlers) urlsByIDHandler(res http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("ClientID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(err)
	}

	urlPairs := handlers.app.GetURLsByUserID(userIDInt)

	if len(urlPairs) == 0 {
		res.WriteHeader(http.StatusNoContent)

	} else {
		var transformedURLPairs []models.URLPair
		var transformedURL models.URLPair

		for _, urlPair := range urlPairs {
			transformedURL.ShortenURL, err = url.JoinPath(handlers.flagConfig.FlagBaseURL, urlPair.ShortenURL)
			if err != nil {
				http.Error(res, "Bad URL path provided", http.StatusInternalServerError)
				log.Println("Failed to join provided URL path with short URL", err)
				return
			}
			transformedURL.OriginalURL = urlPair.OriginalURL
			transformedURLPairs = append(transformedURLPairs, transformedURL)

		}

		resp, err := json.Marshal(transformedURLPairs)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(resp)
	}
}

func (handlers *handlers) deleteURLsHandler(res http.ResponseWriter, req *http.Request) {
	var urls []string
	deleteURLsChannel := make(chan models.URLsClientID, 1)
	defer close(deleteURLsChannel)

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Bad request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(requestBody, &urls)
	if err != nil {
		log.Println("An error occurred while parsing the data", err)
	}

	userID := req.Header.Get("ClientID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Print(err)
	}

	fmt.Println(urls)
	go handlers.app.DeleteURLs(deleteURLsChannel)

	for _, url := range urls {
		var urlsClientID models.URLsClientID
		urlsClientID.URL = url
		urlsClientID.ClientID = userIDInt
		fmt.Println(urlsClientID)
		deleteURLsChannel <- urlsClientID
	}

	res.WriteHeader(http.StatusAccepted)
	res.Write([]byte{})

}
