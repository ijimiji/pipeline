package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ijimiji/pipeline/internal/config"
	"github.com/ijimiji/pipeline/internal/generators/image"
	"github.com/ijimiji/pipeline/internal/instrumentation"
	"github.com/ijimiji/pipeline/internal/processor"
	"github.com/ijimiji/pipeline/internal/services/s3"
	"github.com/ijimiji/pipeline/internal/services/sd"
	"github.com/ijimiji/pipeline/internal/services/sqs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	tracer := instrumentation.NewTracer("sidecar")
	defer tracer.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2114", nil)
	}()

	cfg := config.New[Config](os.Args[len(os.Args)-1])

	sqsClient := sqs.New()
	s3Client := s3.New()
	sd := sd.New(cfg.StableDiffusion)
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
