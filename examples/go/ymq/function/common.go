package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

func sendMessageToQueue(
	ctx context.Context,
	ymqName string,
	message string,
	origin string,
	delay int32,
) (*sqs.SendMessageOutput, error) {
	// Load the SDK's configuration from environment and shared config
	// In the serverless environment, the configuration is loaded from the environment variables
	// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY.
	cfg, err := config.LoadDefaultConfig(
		ctx,
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	// Create an Amazon SQS service client
	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.Region = "ru-central1"
		o.EndpointResolverV2 = &resolverV2{}
	})

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &ymqName,
	}

	urlRes, err := client.GetQueueUrl(ctx, gQInput)

	if err != nil {
		fmt.Printf("Got an error getting the queue URL: %s", err)
		return nil, err
	}

	sMInput := &sqs.SendMessageInput{
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Origin": {
				DataType:    aws.String("String"),
				StringValue: aws.String(origin),
			},
		},
		MessageBody:  aws.String(message),
		QueueUrl:     urlRes.QueueUrl,
		DelaySeconds: delay,
	}
	resp, err := client.SendMessage(ctx, sMInput)
	return resp, nil
}

type resolverV2 struct {
	// you could inject additional application context here as well
}

func (*resolverV2) ResolveEndpoint(_ context.Context, _ sqs.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	u, err := url.Parse("https://message-queue.api.cloud.yandex.net")
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}
