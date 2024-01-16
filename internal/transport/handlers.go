package transport

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/go-chi/chi/v5"
)

func ShortenerHandler(res http.ResponseWriter, req *http.Request) {
	var response string
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Bad request body", http.StatusBadRequest)
		return
	}

	shortenedURL := app.URLMap.ToSortenURL(string(requestBody))

	response, err = url.JoinPath(config.FlagBaseURL, shortenedURL)
	if err != nil {
		log.Println("Ошибка:", err)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(response))
}

func OriginalHandler(res http.ResponseWriter, req *http.Request) {
	idValue := chi.URLParam(req, "id")
	correspondingURL := app.URLMap.ToOriginalURL(idValue)
	res.Header().Set("Location", correspondingURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
