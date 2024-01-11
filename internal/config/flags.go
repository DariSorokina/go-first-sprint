package config

import (
	"flag"
)

var FlagRunAddr string
var FlagBaseURL string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagBaseURL, "b", "http://localhost:8000/", "prefix for trimed URL")
	flag.Parse()
}
