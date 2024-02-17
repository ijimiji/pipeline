package s3

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/ijimiji/pipeline/internal/ptr"
)

func New() *Client {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4566"),
	})

	client := s3.New(sess)

	client.CreateBucket(&s3.CreateBucketInput{
		Bucket: ptr.T("images"),
	})

	client.CreateBucket(&s3.CreateBucketInput{
		Bucket: ptr.T("generation"),
	})

	return &Client{
		originalClient: client,
	}
}

type Client struct {
	originalClient *s3.S3
}

func (c *Client) Get(id string, bucket string) ([]byte, error) {
	out, err := c.originalClient.GetObject(&s3.GetObjectInput{
		Bucket: ptr.T(bucket),
		Key:    ptr.T(id),
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	return io.ReadAll(out.Body)
}

func (c *Client) Put(key string, bucket string, payload []byte) (string, error) {
	_, err := c.originalClient.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(payload),
		Bucket: ptr.T(bucket),
		Key:    ptr.T(key),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}
