package storage

import (
	"context"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestPythonStorageExample(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/python/storage/tf",
		Vars: map[string]interface{}{
			"cloud_id":  os.Getenv("CLOUD_ID"),
			"folder_id": os.Getenv("FOLDER_ID"),
		},
		EnvVars: map[string]string{
			"YC_TOKEN": os.Getenv("YC_TOKEN"),
		},
	})

	terraform.InitAndApply(t, terraformOptions)

	ctx := context.Background()
	bucket := terraform.Output(t, terraformOptions, "bucket_name")
	accessKey := terraform.Output(t, terraformOptions, "access_key")
	secretKey := terraform.Output(t, terraformOptions, "secret_key")

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
		o.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))
	})

	defer func() {
		_, _ = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("uploads/star.png"),
		})
		_, _ = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("thumbnails/star.png"),
		})
	
		terraform.Destroy(t, terraformOptions)
	}()

	// wait for object to be available
	time.Sleep(5 * time.Second)

	filename := "star.png"
	// Open the file for use
	file, err := os.Open(filename)
	assert.NoError(t, err, "Failed to open file %s", filename)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String("uploads/star.png"),
		Body:        file,
		ContentType: aws.String("image/png"),
	})
	assert.NoError(t, err, "Failed to upload file %s to bucket %s", filename, bucket)

	// wait for object to be available
	time.Sleep(5 * time.Second)

	_, err = s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("thumbnails/star.png"),
	})

	assert.NoError(t, err, "Failed to get resized image from bucket %s", bucket)
}

type resolverV2 struct{}

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
