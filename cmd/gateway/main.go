package main

// @title Gateway API
// @version 1.0
// @description Common API for generating images and etc.

// @host localhost:3030
// @BasePath /api/v1

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ijimiji/pipeline/docs"
	"github.com/ijimiji/pipeline/internal/api/http"
	"github.com/ijimiji/pipeline/internal/config"
	"github.com/ijimiji/pipeline/internal/services/core"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.New[Config](os.Args[len(os.Args)-1])

	coreClient := core.New(cfg.Core)

	servant := http.NewServant(cfg.API, coreClient)
	go func() {
		log.Fatal(servant.ListenAndServe())
	}()

	<-done
}
