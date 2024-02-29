package app

import (
	"crypto/md5"
	"fmt"

	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/models"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
)

type App struct {
	storage storage.Database
	log     *logger.Logger
}

func NewApp(storage storage.Database, l *logger.Logger) *App {
	return &App{storage: storage, log: l}
}

func (app *App) ToShortenURL(longURL string, UserID int) (shortURL string, err error) {
	shortURL, err = app.storage.GetShort(longURL)
	if err != nil {
		return
	}
	shortURL = encodeString(longURL)
	app.storage.SetValue(shortURL, longURL, UserID)
	return
}

func (app *App) ToOriginalURL(shortURL string) (longURL string, getOriginalErr error) {
	longURL, getOriginalErr = app.storage.GetOriginal(shortURL)
	return
}

func (app *App) GetURLsByUserID(userID int) (urls []models.URLPair) {
	urls = app.storage.GetURLsByUserID(userID)
	return
}

func (app *App) DeleteURLs(deleteURLsChannel <-chan models.URLsClientID) {
	for urlsClientID := range deleteURLsChannel {
		go app.storage.DeleteURLsWorker(urlsClientID.URL, urlsClientID.ClientID)
	}
}

func (app *App) Ping() (err error) {
	err = app.storage.Ping()
	return
}

func encodeString(data string) string {
	encodedMD5 := md5.Sum([]byte(data))
	encodedMD5Trimed := encodedMD5[:5]
	return fmt.Sprintf("%x", encodedMD5Trimed)

}
