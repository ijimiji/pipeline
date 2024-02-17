package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ijimiji/pipeline/internal/api/grpc"
	"github.com/ijimiji/pipeline/internal/config"
	"github.com/ijimiji/pipeline/internal/managers/image"
	"github.com/ijimiji/pipeline/internal/processor"
	"github.com/ijimiji/pipeline/internal/services/s3"
	"github.com/ijimiji/pipeline/internal/services/sqs"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.New[Config](os.Args[len(os.Args)-1])

	sqsClient := sqs.New()
	s3Client := s3.New()
	imagesManager := image.New(cfg.ImageGeneration, sqsClient, s3Client)
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
