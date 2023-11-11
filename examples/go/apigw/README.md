## Description

This repository contains code of serverless function written in Go that is deployed to Yandex Cloud.
It stands behind the API Gateway and is triggered by HTTP requests. In this example, the function
returns takes a string `name` parametrs from the request body and returns it back.

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

To deploy the infrastructure, run the following command and confirm the action typing `yes`::

```bash
terraform -chdir=./tf apply
```

To test the function, run the following command:

```bash
API_GATEWAY_ENDPOINT=$(terraform -chdir=./tf output -raw api_gateway_endpoint)
curl -XPOST \
  "https://$API_GATEWAY_ENDPOINT/demo" \
  -d '{"name": "test"}' \
  -H "Content-Type: application/json"
```

You should see the following response:

```json
{"message":"Hello, test!"}
```


To destroy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf destroy
```