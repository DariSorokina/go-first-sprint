package config

import (
	"flag"
	"os"
)

type FlagConfig struct {
	FlagRunAddr  string
	FlagBaseURL  string
	FlagLogLevel string
}

func NewFlagConfig() *FlagConfig {
	return &FlagConfig{}
}

func ParseFlags() (flagConfig *FlagConfig) {
	flagConfig = NewFlagConfig()
	flag.StringVar(&flagConfig.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagConfig.FlagBaseURL, "b", "http://localhost:8080/", "prefix for trimed URL")
	flag.StringVar(&flagConfig.FlagLogLevel, "l", "info", "log level")
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
	return
}
