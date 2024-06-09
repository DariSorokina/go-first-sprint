package storage

import (
	"context"
	"errors"

	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/models"
)

var ErrShortURLAlreadyExist = errors.New("corresponding short URL already exists")

type Database interface {
	SetValue(ctx context.Context, shortURL, longURL string, userID int)
	GetShort(ctx context.Context, longURL string) (shortURL string, err error)
	GetOriginal(ctx context.Context, shortURL string) (longURL string, err error)
	GetURLsByUserID(ctx context.Context, userID int) (urls []models.URLPair)
	DeleteURLsWorker(shortURLs []string, userID int)
	Ping(ctx context.Context) error
	Close()
}

func SetStorage(flagConfig *config.FlagConfig, l *logger.Logger) (Database, error) {
	if flagConfig.FlagPostgresqlDSN != "" {
		storage, err := NewPostgresqlDB(flagConfig.FlagPostgresqlDSN, l)
		return storage, err
	}
	storage := NewStorage(flagConfig.FlagFileStoragePath, l)
	return storage, nil
}
