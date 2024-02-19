package instrumentation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"golang.org/x/exp/maps"
)

type SQSCarrier map[string]*sqs.MessageAttributeValue

func (c SQSCarrier) Get(key string) string {
	if c == nil {
		return ""
	}

	attr := c[key]
	if attr == nil {
		return ""
	}

	val := attr.StringValue
	if val == nil {
		return ""
	}

	return *val
}

func (c SQSCarrier) Keys() []string {
	return maps.Keys(c)
}

func (c SQSCarrier) Set(key string, value string) {
	if c == nil {
		return
	}

	c[key] = &sqs.MessageAttributeValue{
		StringValue: aws.String(value),
		DataType:    aws.String("String"),
	}
}
