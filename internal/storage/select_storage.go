// Package storage provides primitives for connecting to data storages.
package storage

import (
	"context"
	"errors"

	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/models"
)

// ErrShortURLAlreadyExist indicates that a corresponding short URL already exists.
var ErrShortURLAlreadyExist = errors.New("corresponding short URL already exists")

// Database is a set of method signatures for data storage.
type Database interface {
	SetValue(ctx context.Context, shortURL, longURL string, userID int)
	GetShort(ctx context.Context, longURL string) (shortURL string, err error)
	GetOriginal(ctx context.Context, shortURL string) (longURL string, err error)
	GetURLsByUserID(ctx context.Context, userID int) (urls []models.URLPair)
	DeleteURLsWorker(shortURLs []string, userID int)
	Ping(ctx context.Context) error
	Close()
}

// SetStorage is a constructor function for data storage object.
func SetStorage(flagConfig *config.FlagConfig, l *logger.Logger) (Database, error) {
	if flagConfig.FlagPostgresqlDSN != "" {
		storage, err := NewPostgresqlDB(flagConfig.FlagPostgresqlDSN, l)
		return storage, err
	}
	storage := NewStorage(flagConfig.FlagFileStoragePath, l)
	return storage, nil
}
