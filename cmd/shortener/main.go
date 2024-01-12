package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/go-chi/chi/v5"
)

var Data = map[string]string{"https://practicum.yandex.ru/": "d41d8cd98f"}
var URLMap = app.URL{Data: Data}

func ShortenerHandler(res http.ResponseWriter, req *http.Request) {
	var response string
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

	if strings.HasSuffix(config.FlagBaseURL, "/") {
		response = config.FlagBaseURL + shortenedURL
	} else {
		response = config.FlagBaseURL + "/" + shortenedURL
	}
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

func run() error {
	fmt.Println("Running server on", config.FlagRunAddr)
	return http.ListenAndServe(config.FlagRunAddr, LinkRouter())
}

func main() {
	config.ParseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}
