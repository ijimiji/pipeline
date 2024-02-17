package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/ijimiji/pipeline/internal/managers/image"
	"github.com/ijimiji/pipeline/internal/models"
	"github.com/ijimiji/pipeline/internal/slices"
	"github.com/ijimiji/pipeline/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(
	config Config,
	imageManager *image.Manager,
) *servant {
	return &servant{
		config:       config,
		imageManager: imageManager,
	}
}

type servant struct {
	proto.UnimplementedCoreServer
	config       Config
	imageManager *image.Manager
}

func (s *servant) Generate(ctx context.Context, req *proto.GenerateRequest) (*proto.GenerateResponse, error) {
	out, err := s.imageManager.Generate(ctx, image.Params{
		Prompt: req.GetPrompt(),
	})

	return &proto.GenerateResponse{
		ID: out.ID,
	}, err
}

func (s *servant) Discard(context.Context, *proto.DiscardRequest) (*proto.DiscardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Discard not implemented")
}

func (s *servant) Status(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	imageGroup, err := s.imageManager.Status(ctx, req.GetID())
	if err != nil {
		return nil, err
	}

	return &proto.StatusResponse{
		ImageGroup: &proto.ImageGroup{
			ID: imageGroup.ID,
			Images: slices.Map(imageGroup.Images, func(image models.Image) *proto.Image {
				return &proto.Image{
					ID:     image.ID,
					Prompt: image.Prompt,
					URL:    image.URL,
					Status: image.GenerationStatus,
				}
			}),
		},
	}, nil
}

func (s *servant) ListenAndServe() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	proto.RegisterCoreServer(server, s)

	return server.Serve(lis)
}
