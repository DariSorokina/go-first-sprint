package app

import (
	"crypto/md5"
	"fmt"

	"github.com/DariSorokina/go-first-sprint.git/internal/models"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
)

type App struct {
	storage storage.Database
}

func NewApp(storage storage.Database) *App {
	return &App{storage: storage}
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

func (app *App) ToOriginalURL(shortURL string) (longURL string) {
	longURL = app.storage.GetOriginal(shortURL)
	return
}

func (app *App) GetURLsByUserID(UserID int) (urls []models.URLPair) {
	urls = app.storage.GetURLsByUserID(UserID)
	return urls
}

func (app *App) Ping() error {
	err := app.storage.Ping()
	return err
}

func encodeString(data string) string {
	encodedMD5 := md5.Sum([]byte(data))
	encodedMD5Trimed := encodedMD5[:5]
	return fmt.Sprintf("%x", encodedMD5Trimed)

}
