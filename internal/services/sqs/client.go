package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ijimiji/pipeline/internal/ptr"
)

func New() *Client {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4566"),
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

func (c *Client) Put(queueName string, message []byte) error {
	_, err := c.originalClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    ptr.T(queueName),
		MessageBody: ptr.T(string(message)),
		// MessageDeduplicationId: ptr.T(uuid.NewString()),
	})

	return err
}

func (c *Client) Recieve(queueName string) ([]byte, error) {
	out, err := c.originalClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            ptr.T(queueName),
		MaxNumberOfMessages: ptr.T[int64](1),
	})
	if err != nil {
		return nil, err
	}
	if len(out.Messages) == 0 {
		return nil, nil
	}

	return []byte(*out.Messages[0].Body), c.Delete(queueName, *out.Messages[0].ReceiptHandle)
}

func (c *Client) Delete(queueName string, id string) error {
	_, err := c.originalClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueName,
		ReceiptHandle: &id,
	})

	return err
}
