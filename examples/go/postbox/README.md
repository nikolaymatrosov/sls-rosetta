## Description

This example demonstrates how send an email from a Yandex Cloud Function using Go via Postbox.

## Prerequisites

* [Terraform](https://www.terraform.io/downloads.html) >= 1.9.7
* [go](https://golang.org/dl/) == 1.21
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
terraform -chdir=./tf apply --auto-approve
```

To test the function, run the following command:

```bash
FUNCTION_ID=$(terraform -chdir=./tf output -raw function_id)
curl "https://functions.yandexcloud.net/$FUNCTION_ID" \
  -H "Content-Type: application/json"
```

You should see the following plain-text response with the message ID:

```
D4UM7H5G7Z1Y.TVC6GXAH7KMR@ingress1-vla
```

To destroy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf destroy
```