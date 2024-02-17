package main

import (
	"context"
	"fmt"
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

	generationProcessor := processor.New(cfg.GenerationProcessor, sqsClient, func(ctx context.Context, req image.GenerateRequest) (image.GenerateResponse, error) {
		var ret image.GenerateResponse
		image, err := sd.Inference(req.Prompt)
		if err != nil {
			return ret, err
		}

		if err := s3Client.Put(req.ID, "generation", image); err != nil {
			return ret, err
		}
		ret.ID = req.ID
		fmt.Println("http://localhost:4566/generation/" + req.ID)

		return ret, nil
	})

	generationContext, generationCancel := context.WithCancel(context.Background())
	go func() {
		generationProcessor.Process(generationContext)
	}()
	<-done
	generationCancel()
}
