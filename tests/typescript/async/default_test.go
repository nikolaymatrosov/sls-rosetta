package async

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/go-resty/resty/v2"
	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTypeScriptAsyncExample(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/typescript/async/tf",
		Vars: map[string]interface{}{
			"cloud_id":  os.Getenv("CLOUD_ID"),
			"folder_id": os.Getenv("FOLDER_ID"),
		},
		EnvVars: map[string]string{
			"YC_TOKEN": os.Getenv("YC_TOKEN"),
		},
	})

	terraform.InitAndApply(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)

	ctx := context.Background()
	ymqName := terraform.Output(t, terraformOptions, "ymq_name")
	accessKey := terraform.Output(t, terraformOptions, "ymq_reader_access_key")
	secretKey := terraform.Output(t, terraformOptions, "ymq_reader_secret_key")

	funcId := terraform.Output(t, terraformOptions, "function_id")

	restClient := resty.New()

	body := map[string]string{
		"name": "test",
	}

	bodyBytes, err := json.Marshal(body)

	httpResp, _ := restClient.R().
		SetHeader("Content-Type", "application/json").
		ForceContentType("application/json").
		SetBody(bodyBytes).
		Post(fmt.Sprintf("https://functions.yandexcloud.net/%s?integration=async", funcId))
	assert.Equal(t, 202, httpResp.StatusCode(), "Status code should be 202")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "https://message-queue.api.cloud.yandex.net",
			SigningRegion: "ru-central1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	ymqClient := sqs.NewFromConfig(cfg)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &ymqName,
	}

	urlRes, err := ymqClient.GetQueueUrl(ctx, gQInput)

	if err != nil {
		t.Errorf("Got an error getting the queue URL: %s", err)
	}

	resp, err := ymqClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            urlRes.QueueUrl,
		MaxNumberOfMessages: 1,
		MessageAttributeNames: []string{
			"All",
		},
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
		WaitTimeSeconds: 20,
	})
	if err != nil {
		fmt.Println("Got an error receiving the message:")
		fmt.Println(err)
		return
	}

	// delete message from queue
	dmInput := &sqs.DeleteMessageInput{
		QueueUrl:      urlRes.QueueUrl,
		ReceiptHandle: resp.Messages[0].ReceiptHandle,
	}
	_, err = ymqClient.DeleteMessage(context.Background(), dmInput)

	print(resp)

	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(*resp.Messages[0].Body), &result)
	if err != nil {
		return
	}
	expected := map[string]interface{}{
		"name": "test",
	}
	if diff := deep.Equal(expected, result); diff != nil {
		t.Error(diff)
	}
}
