package sd

import "context"

type Client interface {
	Inference(context.Context, string) ([]byte, error)
	Close()
}

func New(config Config) Client {
	if config.Stub {
		return newStubClient()
	}

	return newClient()
}
