package apigw

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type result struct {
	Message string `json:"message"`
}

func TestTypescriptApiGatewayExample(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/typescript/apigw/tf",
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

	apiGwEndpoint := terraform.Output(t, terraformOptions, "api_gateway_endpoint")

	client := resty.New()

	var res result
	resp, err := client.R().
		SetBody(`{"name":"test"}`).
		SetHeader("Content-Type", "application/json").
		ForceContentType("application/json").
		SetResult(&res).
		Post(fmt.Sprintf("https://%s/demo", apiGwEndpoint))

	if err != nil {
		t.Errorf("Error sending request to API Gateway: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode(), "Status code should be 200")
	assert.Equal(t, "application/json", resp.Header().Get("content-type"), "Content-Type should be application/json")

	assert.Equal(t, result{Message: "Hello, test!"}, res, "Response body should be {\"message\":\"Hello, test!\"}")
}
