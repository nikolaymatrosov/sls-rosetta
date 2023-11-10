package raw

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type responseBody struct {
	Context map[string]interface{} `json:"context"`
	Request map[string]interface{} `json:"request"`
}

func TestGoRawExample(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/go/raw-function-requests/tf",
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

	var res responseBody

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"message": "Hello, world", "number": 24}`).
		ForceContentType("application/json").
		SetResult(&res).
		Post(fmt.Sprintf("https://functions.yandexcloud.net/%s?integration=raw", functionId))

	if err != nil {
		t.Errorf("Error sending request to API Gateway: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode(), "Status code should be 200, got %d,\n%s", resp.StatusCode(), resp.Body())
	var expected responseBody
	_ = json.Unmarshal([]byte(`{ 
	  "context": {
		"Context": {}
	  },
	  "request": {
		"message": "Hello, world",
		"number": 24
	  }
	}`), &expected)

	assert.Equal(t, expected, res, `Response body doesn't match. Got %s`, resp.Body())
}
