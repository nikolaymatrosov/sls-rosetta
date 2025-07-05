package ydb

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID   int32   `json:"id"`
	Name *string `json:"name"`
}

func TestGoYdbExample(t *testing.T) {
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

	//defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Get outputs from Terraform
	functionUrl := terraform.Output(t, terraformOptions, "function_url")
	ydbEndpoint := terraform.Output(t, terraformOptions, "ydb_endpoint")
	ydbDatabasePath := terraform.Output(t, terraformOptions, "ydb_database_path")

	// Set up the database schema using YDB CLI
	err := setupDatabase(t, ydbEndpoint, ydbDatabasePath)
	if err != nil {
		t.Fatalf("Failed to set up database: %v", err)
	}

	// Wait a bit more for the database to be ready
	//time.Sleep(10 * time.Second)

	// Test the function
	client := resty.New()

	resp, err := client.R().
		Get(functionUrl)

	if err != nil {
		t.Fatalf("Error calling function: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode(), "Status code should be 200")

	// Parse the response
	var user User
	err = json.Unmarshal(resp.Body(), &user)
	if err != nil {
		t.Fatalf("Error parsing response: %v", err)
	}

	// Verify the response
	assert.Equal(t, int32(3), user.ID, "User ID should be 3")
	assert.NotNil(t, user.Name, "User name should not be nil")
	assert.Equal(t, "Charlie", *user.Name, "User name should be Charlie")

	// Test error handling with invalid request
	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"invalid": "data"}).
		Post(functionUrl)

	if err != nil {
		t.Fatalf("Error calling function with invalid request: %v", err)
	}

	// The function should still return 200 for invalid requests since it doesn't process the body
	assert.Equal(t, 200, resp.StatusCode(), "Status code should be 200 for invalid requests")
}

func setupDatabase(t *testing.T, endpoint, databasePath string) error {
	// Check if YDB CLI is available
	_, err := exec.LookPath("ydb")
	if err != nil {
		t.Logf("YDB CLI not found, skipping database setup. Please run setup.sql manually.")
		return nil
	}
	for _, s := range []string{"ddl.sql", "dml.sql"} {

		// Run the setup script
		cmd := exec.Command("ydb",
			"-e", endpoint,
			"-d", databasePath,
			"sql",
			"-f", "../../../examples/go/ydb/"+s,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Failed to run YDB setup script: %v, output: %s", err, string(output))
			// Don't fail the test if YDB CLI setup fails, as it might not be available in CI
			return nil
		}
		t.Logf("Database setup completed: %s", string(output))
	}

	return nil
}
