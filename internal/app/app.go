package app

import (
	"context"
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

func (app *App) ToShortenURL(ctx context.Context, longURL string, userID int) (shortURL string, err error) {
	shortURL, err = app.storage.GetShort(ctx, longURL)
	if err != nil {
		return
	}
	shortURL = encodeString(longURL)
	app.storage.SetValue(ctx, shortURL, longURL, userID)
	return
}

func (app *App) ToOriginalURL(ctx context.Context, shortURL string) (longURL string, err error) {
	longURL, err = app.storage.GetOriginal(ctx, shortURL)
	return
}

func (app *App) GetURLsByUserID(ctx context.Context, userID int) (urls []models.URLPair) {
	urls = app.storage.GetURLsByUserID(ctx, userID)
	return
}

func (app *App) DeleteURLs(deleteURLsChannel <-chan models.URLsClientID) {
	for urlsClientID := range deleteURLsChannel {
		go app.storage.DeleteURLsWorker(urlsClientID.URLs, urlsClientID.ClientID)
	}
}

func (app *App) Ping(ctx context.Context) (err error) {
	err = app.storage.Ping(ctx)
	return
}

func encodeString(data string) string {
	encodedMD5 := md5.Sum([]byte(data))
	encodedMD5Trimed := encodedMD5[:5]
	return fmt.Sprintf("%x", encodedMD5Trimed)

}
