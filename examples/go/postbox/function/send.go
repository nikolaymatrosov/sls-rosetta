package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func sendEmail(ctx context.Context, client *sesv2.Client) (*sesv2.SendEmailOutput, error) {
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
	res, err := client.SendEmail(ctx, input)
	return res, err
}
