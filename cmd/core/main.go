package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ijimiji/pipeline/internal/api/grpc"
	"github.com/ijimiji/pipeline/internal/config"
	"github.com/ijimiji/pipeline/internal/instrumentation"
	"github.com/ijimiji/pipeline/internal/managers/image"
	"github.com/ijimiji/pipeline/internal/processor"
	imagerepository "github.com/ijimiji/pipeline/internal/repositories/image"
	"github.com/ijimiji/pipeline/internal/services/s3"
	"github.com/ijimiji/pipeline/internal/services/sqlite"
	"github.com/ijimiji/pipeline/internal/services/sqs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	tracer := instrumentation.NewTracer("core")
	defer tracer.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	cfg := config.New[Config](os.Args[len(os.Args)-1])

	sqsClient := sqs.New()
	s3Client := s3.New()

	db := sqlite.New()
	imageRepository := imagerepository.New(db)

	imagesManager := image.New(cfg.ImageGeneration, sqsClient, s3Client, imageRepository)
	imageGenerationProcessor := processor.New(cfg.ResultsProcessor, sqsClient, imagesManager.Process)
	resultsContext, resultsCancel := context.WithCancel(context.Background())
	go func() {
		imageGenerationProcessor.Process(resultsContext)
	}()

	servant := grpc.New(cfg.API, imagesManager)
	go func() {
		log.Fatal(servant.ListenAndServe())
	}()

	<-done
	resultsCancel()
}
