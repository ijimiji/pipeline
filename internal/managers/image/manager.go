package image

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ijimiji/pipeline/internal/generators/image"
	"github.com/ijimiji/pipeline/internal/models"
	imagerepository "github.com/ijimiji/pipeline/internal/repositories/image"
	"github.com/ijimiji/pipeline/internal/services/s3"
	"github.com/ijimiji/pipeline/internal/services/sqs"
)

func New(config Config, sqs *sqs.Client, s3 *s3.Client, imageRepository *imagerepository.Repository) *Manager {
	return &Manager{
		config:          config,
		sqs:             sqs,
		s3:              s3,
		imageRepository: imageRepository,
	}
}

type Manager struct {
	config          Config
	sqs             *sqs.Client
	s3              *s3.Client
	imageRepository *imagerepository.Repository
}

type Params struct {
	Prompt string
}

func (m *Manager) Generate(ctx context.Context, params Params) (image.GenerateResponse, error) {
	id := uuid.NewString()

	if err := m.imageRepository.Add(ctx, models.Image{
		Prompt:           params.Prompt,
		ID:               id,
		GenerationStatus: models.GenerationStatusEnqueued,
	}); err != nil {
		return image.GenerateResponse{}, err
	}

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

func (m *Manager) Status(ctx context.Context, id string) (models.ImageGroup, error) {
	var ret models.ImageGroup

	image, err := m.imageRepository.Get(ctx, id)
	if err != nil {
		return ret, err
	}

	ret.ID = id
	ret.Images = append(ret.Images, image)

	return ret, nil
}

func (m *Manager) Process(ctx context.Context, image image.GenerateResponse) (struct{}, error) {
	var ret struct{}

	generatedImage, err := m.s3.Get(image.ID, "generation")
	if err != nil {
		return ret, err
	}

	if err := m.s3.Put(image.ID, "images", generatedImage); err != nil {
		return ret, err
	}

	if err := m.imageRepository.SetURL(ctx, image.ID, "http://localhost:4566/images/"+image.ID); err != nil {
		return ret, err
	}

	return ret, nil
}
