package main

import (
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/transport"
)

func main() {
	config.ParseFlags()
	if err := transport.RunLinkRouter(); err != nil {
		panic(err)
	}
}
