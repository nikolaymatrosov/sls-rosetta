package ydb

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	Rows  []map[string]any `json:"rows"`
	Stats any              `json:"stats"`
}

func TestTypescriptYdbExample(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/typescript/ydb/tf",
		Vars: map[string]interface{}{
			"cloud_id":  os.Getenv("CLOUD_ID"),
			"folder_id": os.Getenv("FOLDER_ID"),
		},
		EnvVars: map[string]string{
			"YC_TOKEN": os.Getenv("YC_TOKEN"),
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	// Migrations run automatically as part of terraform apply
	terraform.InitAndApply(t, terraformOptions)

	functionUrl := terraform.Output(t, terraformOptions, "function_url")

	client := resty.New()

	resp, err := client.R().Get(functionUrl)
	if err != nil {
		t.Fatalf("Error calling function: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode(), "Status code should be 200")

	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		t.Fatalf("Error parsing response: %v", err)
	}

	assert.NotEmpty(t, response.Rows, "Rows should not be empty")
	assert.NotNil(t, response.Stats, "Stats should be present in response")
}
