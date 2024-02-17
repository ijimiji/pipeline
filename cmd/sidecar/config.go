package main

import (
	"github.com/ijimiji/pipeline/internal/processor"
	"github.com/ijimiji/pipeline/internal/services/sd"
)

type Config struct {
	GenerationProcessor processor.Config
	StableDiffusion     sd.Config
}
