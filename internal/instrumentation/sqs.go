package instrumentation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/exp/maps"
)

func NewSQSCarrier(attrs map[string]*sqs.MessageAttributeValue) propagation.TextMapCarrier {
	return &sqsCarrier{
		attrs: attrs,
	}
}

type sqsCarrier struct {
	attrs map[string]*sqs.MessageAttributeValue
}

func (c *sqsCarrier) Get(key string) string {
	if c.attrs == nil {
		return ""
	}

	attr := c.attrs[key]
	if attr == nil {
		return ""
	}

	val := attr.StringValue
	if val == nil {
		return ""
	}

	return *val
}

func (c *sqsCarrier) Keys() []string {
	return maps.Keys(c.attrs)
}

func (c *sqsCarrier) Set(key string, value string) {
	if c.attrs == nil {
		return
	}

	c.attrs[key] = &sqs.MessageAttributeValue{
		StringValue: aws.String(value),
		DataType:    aws.String("String"),
	}
}
