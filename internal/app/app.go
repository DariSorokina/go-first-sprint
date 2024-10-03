// Package app provides the core application logic for managing URLs and user interactions.
package app

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/DariSorokina/go-first-sprint/internal/logger"
	"github.com/DariSorokina/go-first-sprint/internal/models"
	"github.com/DariSorokina/go-first-sprint/internal/storage"
)

// App is a structure representing the application logic.
type App struct {
	storage storage.Database
	log     *logger.Logger
}

// NewApp is a constructor function to create a new App instance.
func NewApp(storage storage.Database, l *logger.Logger) *App {
	return &App{storage: storage, log: l}
}

// ToShortenURL is a method to shorten a long URL and store it in the database.
func (app *App) ToShortenURL(ctx context.Context, longURL string, userID int) (shortURL string, err error) {
	shortURL, err = app.storage.GetShort(ctx, longURL)
	if err != nil {
		return
	}
	shortURL = encodeString(longURL)
	app.storage.SetValue(ctx, shortURL, longURL, userID)
	return
}

// ToOriginalURL is a method to retrieve the original URL from a short URL.
func (app *App) ToOriginalURL(ctx context.Context, shortURL string) (longURL string, err error) {
	longURL, err = app.storage.GetOriginal(ctx, shortURL)
	return
}

// GetURLsByUserID is a method to retrieve URLs associated with a specific user ID.
func (app *App) GetURLsByUserID(ctx context.Context, userID int) (urls []models.URLPair) {
	urls = app.storage.GetURLsByUserID(ctx, userID)
	return
}

// DeleteURLs is a method to handle deletion of URLs based on client requests.
func (app *App) DeleteURLs(deleteURLsChannel <-chan models.URLsClientID) {
	for urlsClientID := range deleteURLsChannel {
		go app.storage.DeleteURLsWorker(urlsClientID.URLs, urlsClientID.ClientID)
	}
}

// Ping is a method to check the database connectivity.
func (app *App) Ping(ctx context.Context) (err error) {
	err = app.storage.Ping(ctx)
	return
}

func encodeString(data string) string {
	encodedMD5 := md5.Sum([]byte(data))
	encodedMD5Trimed := encodedMD5[:5]
	return fmt.Sprintf("%x", encodedMD5Trimed)

}
