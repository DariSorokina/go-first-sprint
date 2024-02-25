package storage

import (
	"errors"

	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/models"
)

var ErrShortURLAlreadyExist = errors.New("corresponding short URL already exists")

type Database interface {
	SetValue(shortURL, longURL string, userID int)
	GetShort(longURL string) (shortURL string, err error)
	GetOriginal(shortURL string) (longURL string)
	GetURLsByUserID(userID int) (urls []models.URLPair)
	Ping() error
	Close()
}

func SetStorage(flagConfig *config.FlagConfig) (storage Database) {
	if flagConfig.FlagPostgresqlDSN != "" {
		storage = NewPostgresqlDB(flagConfig.FlagPostgresqlDSN)
		return
	}
	storage = NewStorage(flagConfig.FlagFileStoragePath)
	return
}
