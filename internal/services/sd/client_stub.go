package sd

import (
	"context"
	_ "embed"
)

func newStubClient() *stubClient {
	return &stubClient{}
}

type stubClient struct{}

func (*stubClient) Close() {
}

//go:embed cat.png
var catImage []byte

func (s *stubClient) Inference(ctx context.Context, prompt string) ([]byte, error) {
	return catImage, nil
}
