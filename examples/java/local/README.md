## Description

This repository contains code of serverless function written in Java that is deployed to Yandex Cloud.
The function is made public and can be triggered by HTTP requests. In this example, the function
takes integer from post body and returns it as a string.

This example shows how to download the function dependencies using Maven and deploy the function using Terraform.

## Prerequisites

* [Terraform](https://www.terraform.io/downloads.html) >= 1.0.0
* [Java](https://www.java.com/en/download/) >= 21
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

To download the function dependencies, run the following command:

```bash
mvn -f ./function/pom.xml dependency:copy-dependencies -DoutputDirectory=.
```

To deploy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf apply
```

To test the function, run the following command:

```bash
FUNCTION_ID=$(terraform -chdir=./tf output -raw function_id)
curl "https://functions.yandexcloud.net/$FUNCTION_ID?integration=raw" \
    --request POST \
    --data '42'
```

You should see the following plain-text response:

```
42
```

To destroy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf destroy
```