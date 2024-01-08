package main

import (
	"io"
	"net/http"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"

	"github.com/gorilla/mux"
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

	vars := mux.Vars(req)
	idValue := vars["id"]

	correspondingURL := URLMap.ToOriginalURL(string(idValue))
	res.Header().Set("Location", correspondingURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", ShortenerHandler).Methods(http.MethodPost)
	router.HandleFunc("/{id}", OriginalHandler).Methods(http.MethodGet)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
