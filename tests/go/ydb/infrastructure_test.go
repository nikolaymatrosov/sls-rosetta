package ydb

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestYdbInfrastructure(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/go/ydb/tf",
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

	// Test that all required outputs are present and not empty
	functionId := terraform.Output(t, terraformOptions, "function_id")
	functionUrl := terraform.Output(t, terraformOptions, "function_url")
	ydbDatabasePath := terraform.Output(t, terraformOptions, "ydb_database_path")
	ydbEndpoint := terraform.Output(t, terraformOptions, "ydb_endpoint")

	// Verify outputs are not empty
	assert.NotEmpty(t, functionId, "Function ID should not be empty")
	assert.NotEmpty(t, functionUrl, "Function URL should not be empty")
	assert.NotEmpty(t, ydbDatabasePath, "YDB database path should not be empty")
	assert.NotEmpty(t, ydbEndpoint, "YDB endpoint should not be empty")

	// Verify URL format
	assert.Contains(t, functionUrl, "https://", "Function URL should be HTTPS")
	assert.Contains(t, functionUrl, "functions.yandexcloud.net", "Function URL should contain Yandex Cloud Functions domain")

	// Verify YDB endpoint format
	assert.Contains(t, ydbEndpoint, "ydb.serverless.yandexcloud.net", "YDB endpoint should contain Yandex Cloud YDB domain")

	// Verify database path format
	assert.Contains(t, ydbDatabasePath, "/ru-central1/", "Database path should contain region")
	assert.Contains(t, ydbDatabasePath, "/", "Database path should contain folder structure")
}

func TestYdbInfrastructureVariables(t *testing.T) {
	// Test that required environment variables are set
	cloudId := os.Getenv("CLOUD_ID")
	folderId := os.Getenv("FOLDER_ID")
	ycToken := os.Getenv("YC_TOKEN")

	// These tests will be skipped if environment variables are not set
	// This allows the test to run in CI environments where these might be set differently
	if cloudId == "" {
		t.Skip("CLOUD_ID environment variable not set")
	}
	if folderId == "" {
		t.Skip("FOLDER_ID environment variable not set")
	}
	if ycToken == "" {
		t.Skip("YC_TOKEN environment variable not set")
	}

	assert.NotEmpty(t, cloudId, "CLOUD_ID should not be empty")
	assert.NotEmpty(t, folderId, "FOLDER_ID should not be empty")
	assert.NotEmpty(t, ycToken, "YC_TOKEN should not be empty")
}
