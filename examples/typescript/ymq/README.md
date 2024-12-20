## Description

This repository contains code of two serverless functions written in Go that is deployed to Yandex Cloud.
First function is triggered by HTTP requests and puts messages to the Yandex Message Queue.
Second function is triggered by Yandex Message Queue and prints messages to the console. It also
puts messages to the Yandex Message Queue for testing purposes.

The `receiver` and `sender` functions intentionally have different style of defining types for arguments and return
values.
It was made to demonstrate that both styles are supported by the Yandex Cloud Functions, and you can choose the one that
you like more.

## Prerequisites

* [Terraform](https://www.terraform.io/downloads.html) >= 1.0.0
* [Node.js](https://nodejs.org/en/download/) >= 18.6.1
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
SEND_FUNC_ID=$(terraform -chdir=./tf output -raw sender_function_id)
curl -XPOST \
  "https://functions.yandexcloud.net/$SEND_FUNC_ID?integration=raw" \
  -d '{"message": "Hello, world", "number": 24}' \
  -H "Content-Type: application/json"
```

To destroy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf destroy
```