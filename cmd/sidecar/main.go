package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ijimiji/pipeline/internal/config"
	"github.com/ijimiji/pipeline/internal/generators/image"
	"github.com/ijimiji/pipeline/internal/processor"
	"github.com/ijimiji/pipeline/internal/services/s3"
	"github.com/ijimiji/pipeline/internal/services/sd"
	"github.com/ijimiji/pipeline/internal/services/sqs"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.New[Config](os.Args[len(os.Args)-1])

	sqsClient := sqs.New()
	s3Client := s3.New()
	sd := sd.New()
	defer sd.Close()

	imageGenerator := image.New(s3Client, sd)

	generationProcessor := processor.New(cfg.GenerationProcessor, sqsClient, imageGenerator.Process)

	generationContext, generationCancel := context.WithCancel(context.Background())
	go func() {
		generationProcessor.Process(generationContext)
	}()
	<-done
	generationCancel()
}
