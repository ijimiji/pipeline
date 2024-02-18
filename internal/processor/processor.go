package processor

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ijimiji/pipeline/internal/services/sqs"
)

type processFunc[Request any, Response any] func(ctx context.Context, request Request) (Response, error)

func New[Request any, Response any](config Config, sqsClient *sqs.Client, pFunc processFunc[Request, Response]) Processor[Request, Response] {
	return Processor[Request, Response]{
		pFunc:  pFunc,
		sqs:    sqsClient,
		config: config,
	}
}

type Processor[Request any, Response any] struct {
	config Config
	sqs    *sqs.Client
	pFunc  processFunc[Request, Response]
}

func (p *Processor[Request, Response]) Process(ctx context.Context) error {
	for {
		time.Sleep(time.Second * 4)
		select {
		case <-ctx.Done():
			return nil
		default:
			message, err := p.sqs.Recieve(p.config.InputQueue)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if len(message) == 0 {
				continue
			}

			var req Request
			if err := json.Unmarshal(message, &req); err != nil {
				slog.Error(err.Error())
				continue
			}

			resp, err := p.pFunc(ctx, req)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if len(p.config.OutputQueue) == 0 {
				continue
			}

			payload, err := json.Marshal(resp)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if err := p.sqs.Put(p.config.OutputQueue, payload); err != nil {
				slog.Error(err.Error())
			}
		}
	}
}
