package processor

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ijimiji/pipeline/internal/services/sqs"
	"go.opentelemetry.io/otel"
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
		time.Sleep(time.Second * 1)
		select {
		case <-ctx.Done():
			return nil
		default:
			func(ctx context.Context) {
				message, carrier, err := p.sqs.Recieve(ctx, p.config.InputQueue)
				if err != nil {
					slog.Error(err.Error())
					return
				}

				if len(message) == 0 {
					return
				}

				tracer := otel.Tracer("processor")
				ctx, span := tracer.Start(otel.GetTextMapPropagator().Extract(ctx, carrier), "processor")
				defer span.End()

				var req Request
				if err := json.Unmarshal(message, &req); err != nil {
					slog.Error(err.Error())
					return
				}

				resp, err := p.pFunc(ctx, req)
				if err != nil {
					slog.Error(err.Error())
					return
				}

				if len(p.config.OutputQueue) == 0 {
					return
				}

				payload, err := json.Marshal(resp)
				if err != nil {
					slog.Error(err.Error())
					return
				}

				if err := p.sqs.Put(ctx, p.config.OutputQueue, payload); err != nil {
					slog.Error(err.Error())
				}
			}(context.Background())
		}
	}
}
