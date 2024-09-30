package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

// Handler handles an object storage event.
// It creates a new S3 client, retrieves the object involved in the event, and returns a response.
func Handler(ctx context.Context, event *ObjectStorageEvent) (*ObjectStorageResponse, error) {
	// Load the AWS configuration with the custom endpoint resolver.
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithDefaultRegion("ru-central1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new S3 client.
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "ru-central1"
		o.EndpointResolverV2 = &resolverV2{}
	})
	// Initialize a WaitGroup to manage the goroutine.
	wg := sync.WaitGroup{}
	// Add a task to the WaitGroup.
	wg.Add(len(event.Messages))

	for _, message := range event.Messages {
		// Get the object involved in the event.
		object, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(message.Details.BucketID),
			Key:    aws.String(message.Details.ObjectID),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get object: %w", err)
		}

		// Print the size of the object to stdout.
		fmt.Println("Object size:", object.ContentLength)
		thumbnailKey := "thumbnail/" + strings.TrimPrefix(message.Details.ObjectID, "uploads/")

		pipeReader, pipeWriter := io.Pipe()

		// Start a new goroutine to handle the object storage operation.
		go func() {
			// Attempt to put the object into the bucket.
			_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
				// The name of the bucket to put the object into.
				Bucket: aws.String(message.Details.BucketID),

				// The key to store the object under, prefixed with "thumbnail-".
				Key: aws.String(thumbnailKey),

				// The data to store in the object.
				Body: pipeReader,

				// The MIME type of the object.
				ContentType: object.ContentType,
			})

			// If an error occurred while putting the object, panic.
			if err != nil {
				panic(fmt.Errorf("failed to upload object: %w", err))
			}

			// Signal to the WaitGroup that the task is done.
			wg.Done()
		}()

		// Create a thumbnail of the object.
		defer object.Body.Close()
		err = Thumbnail(object.Body, pipeWriter)
		if err != nil {
			return nil, err
		}
	}
	// Wait for all goroutins to finish.
	wg.Wait()

	// Return a successful response.
	return &ObjectStorageResponse{
		StatusCode: 200,
	}, nil
}

type resolverV2 struct {
	// you could inject additional application context here as well
}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	u, err := url.Parse("https://storage.yandexcloud.net")
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	u.Path += "/" + *params.Bucket
	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}
