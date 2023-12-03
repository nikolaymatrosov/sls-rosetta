package hello

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTypescriptHelloExample(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/typescript/hello/tf",
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

	functionId := terraform.Output(t, terraformOptions, "function_id")

	client := resty.New()

	resp, err := client.R().
		SetQueryParam("name", "test").
		Get(fmt.Sprintf("https://functions.yandexcloud.net/%s", functionId))

	if err != nil {
		t.Errorf("Error sending request to API Gateway: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode(), "Status code should be 200")

	assert.Equal(t, []byte("Hello, test!"), resp.Body(), `Response body should be "Hello, test!"`)
}
