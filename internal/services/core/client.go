package core

import (
	"github.com/ijimiji/pipeline/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	proto.CoreClient
}

func New(config Config) Client {
	conn, err := grpc.Dial(config.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
