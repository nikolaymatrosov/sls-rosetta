package documentdb

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type items []map[string]interface{}

func TestTypescriptDocumentDBExample(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../../examples/typescript/documentdb/tf",
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

	functionId := terraform.Output(t, terraformOptions, "function_id")
	randomId := random.Random(0, 1000)
	client := resty.New()

	resp, err := fetch(client, functionId, randomId)

	if err != nil {
		t.Errorf("Error sending request to function: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode(), "Status code should be 200")

	assert.Equal(t, []byte("[]"), resp.Body(), `Response body should be "[]"`)

	randomValue := uuid.New()

	req := map[string]interface{}{
		"id":    float64(randomId),
		"key":   "test",
		"value": randomValue.String(),
	}
	resp, err = client.R().
		SetBody(req).
		ForceContentType("application/json").
		Post(fmt.Sprintf("https://functions.yandexcloud.net/%s", functionId))

	if err != nil {
		t.Errorf("Error sending request to API Gateway: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode(), "Status code should be 200")

	resp, err = fetch(client, functionId, randomId)

	if err != nil {
		t.Errorf("Error sending request to function: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode(), "Status code should be 200")

	expected := items{
		req,
	}

	var result items
	err = json.Unmarshal(resp.Body(), &result)

	if diff := deep.Equal(expected, result); diff != nil {
		t.Error(diff)
	}
}

func fetch(client *resty.Client, functionId string, id int) (*resty.Response, error) {
	return client.R().
		SetQueryParam("id", strconv.Itoa(id)).
		Get(fmt.Sprintf("https://functions.yandexcloud.net/%s", functionId))

}
