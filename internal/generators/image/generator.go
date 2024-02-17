package image

import (
	"context"
	"fmt"

	"github.com/ijimiji/pipeline/internal/services/s3"
	"github.com/ijimiji/pipeline/internal/services/sd"
)

func New(s3Client *s3.Client, sd sd.Client) *Generator {
	return &Generator{
		s3Client: s3Client,
		sd:       sd,
	}
}

type Generator struct {
	s3Client *s3.Client
	sd       sd.Client
}

func (g *Generator) Process(ctx context.Context, req GenerateRequest) (GenerateResponse, error) {
	var ret GenerateResponse
	image, err := g.sd.Inference(ctx, req.Prompt)
	if err != nil {
		return ret, err
	}

	if err := g.s3Client.Put(req.ID, "generation", image); err != nil {
		return ret, err
	}
	ret.ID = req.ID
	fmt.Println("http://localhost:4566/generation/" + req.ID)

	return ret, nil
}
