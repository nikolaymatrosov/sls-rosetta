## Description

This repository contains code of serverless function written in Typescript that is deployed to Yandex Cloud.
The function is made public and can be triggered by HTTP requests. In this example, the function
takes name from query parameters and returns it back.

## Prerequisites

* [Terraform](https://www.terraform.io/downloads.html) >= 1.0.0
* [Node.js](https://nodejs.org/en/download/) >= 18.16.1
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
FUNCTION_ID=$(terraform -chdir=./tf output -raw function_id)
curl "https://functions.yandexcloud.net/$FUNCTION_ID?name=test" \
  -H "Content-Type: application/json"
```

You should see the following plain-text response:

```
Hello, test!
```

To destroy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf destroy
```