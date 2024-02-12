package storage

import "github.com/DariSorokina/go-first-sprint.git/internal/config"

type Database interface {
	SetValue(shortURL, longURL string)
	GetShort(longURL string) (shortURL string)
	GetOriginal(shortURL string) (longURL string)
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
