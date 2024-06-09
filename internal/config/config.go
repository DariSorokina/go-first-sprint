// Package config provides functionality for managing configuration flags and environment variables.
package config

import (
	"flag"
	"os"
)

// FlagConfig is a structure containing configuration flags for the server.
type FlagConfig struct {
	FlagRunAddr         string
	FlagBaseURL         string
	FlagLogLevel        string
	FlagFileStoragePath string
	FlagPostgresqlDSN   string
}

// NewFlagConfig is a constructor function to create a new FlagConfig instance.
func NewFlagConfig() *FlagConfig {
	return &FlagConfig{}
}

// "/tmp/short-url-db.json"

// ParseFlags is a function to parse command-line flags and environment variables into FlagConfig.
func ParseFlags() (flagConfig *FlagConfig) {
	flagConfig = NewFlagConfig()
	flag.StringVar(&flagConfig.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagConfig.FlagBaseURL, "b", "http://localhost:8080/", "prefix for trimed URL")
	flag.StringVar(&flagConfig.FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&flagConfig.FlagFileStoragePath, "f", "/Users/dariasorokina/Desktop/yp_golang/go-first-sprint/internal/storage/short-url-db.json", "file storage path")
	flag.StringVar(&flagConfig.FlagPostgresqlDSN, "d", "", "postgreSQL DSN")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		flagConfig.FlagRunAddr = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		flagConfig.FlagBaseURL = envBaseURL
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagConfig.FlagLogLevel = envLogLevel
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		flagConfig.FlagFileStoragePath = envFileStoragePath
	}
	if envPostgresqlDSN := os.Getenv("DATABASE_DSN"); envPostgresqlDSN != "" {
		flagConfig.FlagPostgresqlDSN = envPostgresqlDSN
	}
	return
}
