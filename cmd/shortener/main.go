package main

import (
	"io"
	"log"
	"net/http"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/go-chi/chi/v5"
)

var Data = map[string]string{"https://practicum.yandex.ru/": "d41d8cd98f"}
var URLMap = app.URL{Data: Data}

func ShortenerHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Bad request body", http.StatusBadRequest)
		return
	}

	shortenedURL := URLMap.ToSortenURL(string(requestBody))

	response := "http://localhost:8080/" + shortenedURL
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(response))
}

func OriginalHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	idValue := chi.URLParam(req, "id")
	correspondingURL := URLMap.ToOriginalURL(idValue)
	res.Header().Set("Location", correspondingURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func LinkRouter() chi.Router {
	router := chi.NewRouter()
	router.HandleFunc("/", ShortenerHandler)
	router.HandleFunc("/{id}", OriginalHandler)
	return router
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", LinkRouter()))
}
