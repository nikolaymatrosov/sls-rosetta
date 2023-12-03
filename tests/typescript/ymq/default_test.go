package ymq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/go-resty/resty/v2"
	"github.com/go-test/deep"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

type customResolver2 struct {
}

func (c customResolver2) ResolveEndpoint(ctx context.Context, params sqs.EndpointParameters) (
	smithyendpoints.Endpoint,
	error,
) {
	return smithyendpoints.Endpoint{
		URI: url.URL{
			Scheme: "https",
			Host:   "message-queue.api.cloud.yandex.net",
		},
	}, nil
}

func TestTypescriptYmqExample(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/typescript/ymq/tf",
		Vars: map[string]interface{}{
			"cloud_id":  os.Getenv("CLOUD_ID"),
			"folder_id": os.Getenv("FOLDER_ID"),
		},
		EnvVars: map[string]string{
			"YC_TOKEN": os.Getenv("YC_TOKEN"),
		},
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	ctx := context.Background()
	ymqName := terraform.Output(t, terraformOptions, "ymq_name")
	accessKey := terraform.Output(t, terraformOptions, "ymq_reader_access_key")
	secretKey := terraform.Output(t, terraformOptions, "ymq_reader_secret_key")

	senderFuncId := terraform.Output(t, terraformOptions, "sender_function_id")

	restClient := resty.New()

	_, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		ForceContentType("application/json").
		Get(fmt.Sprintf("https://functions.yandexcloud.net/%s?integration=raw", senderFuncId))

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
		config.WithDefaultRegion("ru-central1"),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	ymqClient := sqs.NewFromConfig(cfg,
		sqs.WithEndpointResolverV2(customResolver2{}),
	)

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

	print(resp)
	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(*resp.Messages[0].Body), &result)
	if err != nil {
		return
	}
	expected := map[string]interface{}{
		"result": "success",
		"name":   "test",
	}
	if diff := deep.Equal(expected, result); diff != nil {
		t.Error(diff)
	}

}
