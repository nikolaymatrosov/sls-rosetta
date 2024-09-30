## Description

This repository contains code of serverless function written in Go that is deployed to Yandex Cloud.
The function is made public and can be triggered by HTTP requests.

Then we will try to examine what effect the concurrency has on the function execution time, rps and memory consumption.

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
FUNCTION_ID=$(terraform -chdir=./tf output -raw simple)
curl -XPOST \
  "https://functions.yandexcloud.net/$FUNCTION_ID" \
  -H "X-Request-Id: `uuidgen`" \
  -d '{"name": "John"}' \
  -H "Content-Type: application/json"
```

To destroy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf destroy
```


```bash
FUNCTION_ID=$(terraform -chdir=./tf output -raw ydb)
curl -XPOST \
  "https://functions.yandexcloud.net/$FUNCTION_ID" \
  -H "X-Request-Id: `uuidgen`" \
  -d '{"name": "John"}' \
  -H "Content-Type: application/json"
```