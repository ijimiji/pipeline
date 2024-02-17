package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/ijimiji/pipeline/internal/managers/image"
	"github.com/ijimiji/pipeline/proto"
	"google.golang.org/grpc"
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
	_, err := s.imageManager.Generate(ctx, image.Params{
		Prompt: "very fat bober",
	})
	return &proto.GenerateResponse{}, err
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
