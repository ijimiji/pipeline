package main

import (
	"github.com/ijimiji/pipeline/internal/api/grpc"
	"github.com/ijimiji/pipeline/internal/managers/image"
	"github.com/ijimiji/pipeline/internal/processor"
)

type Config struct {
	API              grpc.Config
	ResultsProcessor processor.Config
	ImageGeneration  image.Config
}
