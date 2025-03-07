package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func AwsHandler(ctx context.Context) ([]byte, error) {
	// Create an SES session.
	client := sesv2.New(sesv2.Options{
		Region:             "ru-central1",
		EndpointResolverV2: &resolverV2{},
		// The following options are useful for debugging
		//ClientLogMode:      aws.LogRequestWithBody | aws.LogResponseWithBody,
		//Logger: logging.NewStandardLogger(
		//	os.Stdout,
		//),
	})

	res, err := sendEmail(ctx, client)

	if err != nil {
		return nil, err
	}

	// return response
	resp := map[string]interface{}{
		"result":    "success",
		"messageId": res.MessageId,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return respBytes, nil

}
