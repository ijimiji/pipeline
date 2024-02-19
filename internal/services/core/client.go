package core

import (
	"github.com/ijimiji/pipeline/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	proto.CoreClient
}

func New(config Config) Client {
	conn, err := grpc.Dial(
		config.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}

	return &client{
		CoreClient: proto.NewCoreClient(conn),
	}
}

type client struct {
	proto.CoreClient
}
