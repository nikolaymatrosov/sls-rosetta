## Description

This repository contains code of serverless function written in Go that is deployed to Yandex Cloud.
The function is made public and can be triggered by HTTP requests. In this example, the function
is invoked asynchronously. The result of the function execution is stored in the Yandex Cloud Message Queue.

## Prerequisites

* [Terraform](https://www.terraform.io/downloads.html) >= 1.0
* [Go](https://golang.org/doc/install) >= 1.21
* [Yandex Cloud CLI](https://cloud.yandex.ru/docs/cli/quickstart)
* [curl](https://curl.se/download.html)

## Usage with Terraform deploy

To initialize Terraform, run the following command:

```bash
terraform -chdir=./tf init
```

To set the environment variables, run the following command:

```bash
export TF_VAR_cloud_id=b1g***
export TF_VAR_folder_id=b1g***
export YC_TOKEN=`yc iam create-token`
```

To deploy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf apply
```

To test the function, run the following command:

```bash
FUNCTION_ID=$(terraform -chdir=./tf output -raw function_id)
curl "https://functions.yandexcloud.net/$FUNCTION_ID?integration=async" \
  -H "Content-Type: application/json" \
  -d '{"name": "test"}' \
  -X POST
```

In the successful case, you should see the following message in the queue:

```bash
QUEUE_URL=$(terraform -chdir=./tf output -raw ymq_id)
AWS_ACCESS_KEY_ID=$(terraform -chdir=./tf output -raw ymq_reader_access_key)
AWS_SECRET_ACCESS_KEY=$(terraform -chdir=./tf output -raw ymq_reader_secret_key)
aws sqs receive-message --queue-url $QUEUE_URL --endpoint https://message-queue.api.cloud.yandex.net
```

You should see response like this:
```json
{
    "Messages": [
        {
            "MessageId": "d387065-774b1206-85bf6f17-73b38758",
            "ReceiptHandle": "EAEgw8PAq80xKAI",
            "MD5OfBody": "e2b1c523c4703d2780f0f0a0fddbb905",
            "Body": "{\"name\":\"test\",\"result\":\"success\"}",
            "Attributes": {
                "ApproximateFirstReceiveTimestamp": "1704387944899",
                "ApproximateReceiveCount": "1",
                "SentTimestamp": "1704387941053",
                "SenderId": "aje2upl6d1anmqppsamg@as"
            }
        }
    ]
}
```

To destroy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf destroy
```