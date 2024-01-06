package main

import (
	"io"
	"log"
	"net/http"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"

	"github.com/gorilla/mux"
)

var URLMap = app.URL{}

func apiShortener(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		res.Write([]byte("Плохое тело запроса"))
		return
	}

	shortenedURL := URLMap.ToSortenURL(string(requestBody))

	response := "http://localhost:8080/" + shortenedURL
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(response))
}

func apiOriginal(res http.ResponseWriter, req *http.Request) {
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

func NotFoundHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Bad request", http.StatusBadRequest)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", apiShortener).Methods(http.MethodPost)
	router.HandleFunc("/{id}", apiOriginal).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
