package image

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/ijimiji/pipeline/internal/generators/image"
	"github.com/ijimiji/pipeline/internal/services/s3"
	"github.com/ijimiji/pipeline/internal/services/sqs"
)

func New(config Config, sqs *sqs.Client, s3 *s3.Client) *Manager {
	return &Manager{
		config: config,
		sqs:    sqs,
		s3:     s3,
	}
}

type Manager struct {
	config Config
	sqs    *sqs.Client
	s3     *s3.Client
}

type Params struct {
	Prompt string
}

func (m *Manager) Generate(ctx context.Context, params Params) (image.GenerateResponse, error) {
	id := uuid.NewString()

	marshaled, err := json.Marshal(image.GenerateRequest{
		ID:     id,
		Prompt: params.Prompt,
	})
	if err != nil {
		return image.GenerateResponse{}, err
	}

	if err := m.sqs.Put(m.config.GenerationQueue, marshaled); err != nil {
		return image.GenerateResponse{}, err
	}

	return image.GenerateResponse{
		ID: id,
	}, nil
}

func (m *Manager) Process(ctx context.Context, image image.GenerateResponse) (struct{}, error) {
	var ret struct{}

	generatedImage, err := m.s3.Get(image.ID, "generation")
	if err != nil {
		return ret, err
	}

	key, err := m.s3.Put(image.ID, "images", generatedImage)
	if err != nil {
		return ret, err
	}

	fmt.Println(key)

	return ret, nil
}
