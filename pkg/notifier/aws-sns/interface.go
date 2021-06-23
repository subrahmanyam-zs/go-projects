package awssns

import "github.com/aws/aws-sdk-go/service/sns"

type AWS interface {
	Subscribe(input *sns.SubscribeInput) (*sns.SubscribeOutput, error)
	Publish(input *sns.PublishInput) (*sns.PublishOutput, error)
}
