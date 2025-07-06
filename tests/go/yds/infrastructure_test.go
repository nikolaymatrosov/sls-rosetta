package yds

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestYdsInfrastructure(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/go/yds/tf",
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

	outputs := []string{
		"producer_function_url",
		"consumer_function_id",
		"yds_topic_id",
		"yds_topic_name",
		"yds_database_id",
		"yds_database_path",
		"trigger_id",
	}

	for _, output := range outputs {
		val := terraform.Output(t, terraformOptions, output)
		assert.NotEmpty(t, val, output+" should not be empty")
	}
}
