package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/DariSorokina/go-first-sprint.git/internal/app"
	"github.com/DariSorokina/go-first-sprint.git/internal/config"
	"github.com/DariSorokina/go-first-sprint.git/internal/logger"
	"github.com/DariSorokina/go-first-sprint.git/internal/server"
	"github.com/DariSorokina/go-first-sprint.git/internal/storage"
)

func shortner() {
	flagConfig := config.ParseFlags()

	var l *logger.Logger
	var err error
	if l, err = logger.CreateLogger(flagConfig.FlagLogLevel); err != nil {
		log.Fatal("Failed to create logger:", err)
	}

	storage, err := storage.SetStorage(flagConfig, l)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	app := app.NewApp(storage, l)
	serv := server.NewServer(app, flagConfig, l)

	if err := server.Run(serv); err != nil {
		panic(err)
	}
}

func memProfile() {
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func main() {
	go shortner()
	time.Sleep(30 * time.Second)

	memProfile()
}
