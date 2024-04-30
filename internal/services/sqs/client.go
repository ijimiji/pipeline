package sqs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ijimiji/pipeline/internal/instrumentation"
	"github.com/ijimiji/pipeline/internal/ptr"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func New() *Client {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4566"),
		HTTPClient: &http.Client{
			Transport: otelhttp.NewTransport(
				http.DefaultTransport,
				otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
					return fmt.Sprintf("sqs HTTP %s", r.Method)
				}),
			),
		},
	})

	client := sqs.New(sess)

	client.CreateQueue(&sqs.CreateQueueInput{
		QueueName: ptr.T("generation"),
	})

	client.CreateQueue(&sqs.CreateQueueInput{
		QueueName: ptr.T("results"),
	})

	return &Client{
		originalClient: client,
	}
}

type Client struct {
	originalClient *sqs.SQS
}

func (c *Client) Put(ctx context.Context, queueName string, message []byte) error {
	attrs := map[string]*sqs.MessageAttributeValue{}
	otel.GetTextMapPropagator().Inject(ctx, instrumentation.SQSCarrier(attrs))

	_, err := c.originalClient.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		QueueUrl:          ptr.T(queueName),
		MessageBody:       ptr.T(string(message)),
		MessageAttributes: attrs,
	})

	return err
}

func (c *Client) Recieve(ctx context.Context, queueName string) ([]byte, propagation.TextMapCarrier, error) {
	out, err := c.originalClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              ptr.T(queueName),
		MaxNumberOfMessages:   ptr.T[int64](1),
		MessageAttributeNames: []*string{aws.String("All")},
	})
	if err != nil {
		return nil, nil, err
	}
	if len(out.Messages) == 0 {
		return nil, nil, nil
	}
	message := out.Messages[0]
	carrier := instrumentation.SQSCarrier(message.MessageAttributes)

	return []byte(*message.Body), carrier, c.Delete(ctx, queueName, *out.Messages[0].ReceiptHandle)
}

func (c *Client) Delete(ctx context.Context, queueName string, id string) error {
	_, err := c.originalClient.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueName,
		ReceiptHandle: &id,
	})

	return err
}
