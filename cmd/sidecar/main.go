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
	"github.com/ijimiji/pipeline/internal/services/sqs"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.New[Config](os.Args[len(os.Args)-1])

	sqsClient := sqs.New()
	generationProcessor := processor.New(cfg.GenerationProcessor, sqsClient, func(ctx context.Context, req image.GenerateRequest) (image.GenerateResponse, error) {
		fmt.Println(req.Prompt)
		return image.GenerateResponse{}, nil
	})

	generationContext, generationCancel := context.WithCancel(context.Background())
	go func() {
		generationProcessor.Process(generationContext)
	}()
	<-done
	generationCancel()
}
