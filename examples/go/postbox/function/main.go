package main

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

const (
	// Sender address must be verified with Amazon SES.
	Sender = "noreply@yourdomain.com"

	// Recipient address.
	Recipient = "receiver@domain.com"

	// Subject line for the email.
	Subject = "Yandex Cloud Postbox Test via AWS SDK for Go"

	// HtmlBody is the body for the email.
	HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
		"<a href='https://yandex.cloud/ru/docs/postbox/quickstart'>Yandex Cloud Postbox</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	// TextBody is the email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Yandex Cloud Postbox using the AWS SDK for Go."

	// CharSet The character encoding for the email.
	CharSet = "UTF-8"
)

func Handler(_ context.Context) ([]byte, error) {
	// Create an SES session.
	client := sesv2.New(sesv2.Options{
		Region:             "ru-central1",
		EndpointResolverV2: &resolverV2{},
		//ClientLogMode:      aws.LogRequestWithBody | aws.LogResponseWithBody,
		//Logger: logging.NewStandardLogger(
		//	os.Stdout,
		//),
	})

	// Assemble the email.
	input := &sesv2.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{Recipient},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(Subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(HtmlBody),
					},
					Text: &types.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(TextBody),
					},
				},
			},
		},
		FromEmailAddress: aws.String(Sender),
	}

	// Attempt to send the email.
	res, err := client.SendEmail(context.Background(), input)

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

type resolverV2 struct{}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params sesv2.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	u, err := url.Parse("https://postbox.cloud.yandex.net")
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}
