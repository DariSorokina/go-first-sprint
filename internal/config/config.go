package config

import (
	"flag"
	"os"
)

var FlagRunAddr string
var FlagBaseURL string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagBaseURL, "b", "http://localhost:8080/", "prefix for trimed URL")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		FlagBaseURL = envBaseURL
	}
}
