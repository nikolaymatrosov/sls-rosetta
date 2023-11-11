## Description

This repository contains code of two serverless functions written in Go that is deployed to Yandex Cloud.
First function is triggered by HTTP requests and puts messages to the Yandex Data Stream.
Second function is triggered by Yandex Data Stream and prints messages to the console. It also
puts messages to the Yandex Message Queue for testing purposes.

## Prerequisites

* [Terraform](https://www.terraform.io/downloads.html) >= 0.14.0
* [Go](https://golang.org/doc/install) >= 1.19
* [Yandex Cloud CLI](https://cloud.yandex.ru/docs/cli/quickstart)
* [curl](https://curl.se/download.html)

## Usage with Terraform deploy

To initialize Terraform, run the following command:

```bash
terraform -chdir=./tf init
```

To set the environment variables, run the following command:

```bash
export TF_VAR_cloud_id=b1grv4795cleikron6a9
export TF_VAR_folder_id=b1g8l63q7v4dqkl3bnkj
export YC_TOKEN=`yc iam create-token`
```

To deploy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf apply
```

To test the function, run the following command:

```bash
SEND_FUNC_ID=$(terraform -chdir=./tf output -raw send_function_id)
curl -XPOST \
  "https://functions.yandexcloud.net/$SEND_FUNC_ID?integration=raw" \
  -d '{"message": "Hello, world", "number": 24}' \
  -H "Content-Type: application/json"
```

To destroy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf destroy
```