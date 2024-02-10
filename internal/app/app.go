package app

import (
	"crypto/md5"
	"fmt"

	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
)

type App struct {
	storage *storage.Storage
}

func NewApp(storage *storage.Storage) *App {
	return &App{storage: storage}
}

func (app *App) ToShortenURL(longURL string) (shortURL string) {
	shortURL = app.storage.GetShort(longURL)
	if shortURL == "" {
		shortURL = encodeString(longURL)
		app.storage.SetValue(shortURL, longURL)
	}
	return
}

func (app *App) ToOriginalURL(shortURL string) (longURL string) {
	longURL = app.storage.GetOriginal(shortURL)
	return
}

// а точно ли нужно метод прокидывать так, или есть альтернативы?
func (app *App) PingPostgresql() error {
	err := app.storage.PingPostgresql()
	return err
}

func encodeString(data string) string {
	encodedMD5 := md5.Sum([]byte(data))
	encodedMD5Trimed := encodedMD5[:5]
	return fmt.Sprintf("%x", encodedMD5Trimed)

}
