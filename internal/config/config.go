package config

import (
	"flag"
	"os"
)

type FlagConfig struct {
	FlagRunAddr         string
	FlagBaseURL         string
	FlagLogLevel        string
	FlagFileStoragePath string
}

func NewFlagConfig() *FlagConfig {
	return &FlagConfig{}
}

func ParseFlags() (flagConfig *FlagConfig) {
	flagConfig = NewFlagConfig()
	flag.StringVar(&flagConfig.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagConfig.FlagBaseURL, "b", "http://localhost:8080/", "prefix for trimed URL")
	flag.StringVar(&flagConfig.FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&flagConfig.FlagFileStoragePath, "f", "/Users/dariasorokina/Desktop/yp_golang/go-first-sprint/cmd/shortener/short-url-db.json", "file storage path")
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
	return
}
