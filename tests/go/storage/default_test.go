package ymq

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestGoStorageExample(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/go/storage/tf",
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
	bucket := terraform.Output(t, terraformOptions, "bucket")
	bucketForFunction := terraform.Output(t, terraformOptions, "bucket_for_function")
	accessKey := terraform.Output(t, terraformOptions, "sa_storage_editor_access_key")
	secretKey := terraform.Output(t, terraformOptions, "sa_storage_editor_secret_key")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "https://storage.yandexcloud.net",
			SigningRegion: "ru-central1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)

	s3Client := s3.NewFromConfig(cfg)

	defer func() {
		_, _ = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("uploads/star.png"),
		})
		_, _ = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("thumbnails/star.png"),
		})
		_, _ = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucketForFunction),
			Key:    aws.String("function.zip"),
		})

		terraform.Destroy(t, terraformOptions)
	}()

	// wait for object to be available
	time.Sleep(30 * time.Second)

	filename := "star.png"
	// Open the file for use
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Got error opening file:")
		fmt.Println(err)
		return
	}

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String("uploads/star.png"),
		Body:        file,
		ContentType: aws.String("image/png"),
	})

	if err != nil {
		fmt.Println("Got an error receiving the message:")
		fmt.Println(err)
		return
	}

	// wait for object to be available
	time.Sleep(5 * time.Second)

	_, err = s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("thumbnail/star.png"),
	})

	if err != nil {
		fmt.Println("Got an error fetching resized image:")
		fmt.Println(err)
		t.Fail()
		return
	}
}
