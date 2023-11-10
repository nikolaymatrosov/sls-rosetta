## Description

This repository contains code of serverless function written in Go that is deployed to Yandex Cloud.
It stands behind the API Gateway and is triggered by HTTP requests. In this example, the function
returns takes a string `name` parametrs from the request body and returns it back.

## Prerequisites

* [Terraform](https://www.terraform.io/downloads.html) >= 0.14.0
* [Go](https://golang.org/doc/install) >= 1.19
* [Yandex Cloud CLI](https://cloud.yandex.ru/docs/cli/quickstart)
* [curl](https://curl.se/download.html)

## Usage

To initialize Terraform, run the following command:

```bash
make init
```

To set the environment variables, run the following command:

```bash
export TF_VAR_cloud_id=b1grv4795cleikron6a9
export TF_VAR_folder_id=b1g8l63q7v4dqkl3bnkj
export YC_TOKEN=`yc iam create-token`
```

To deploy the infrastructure, run the following command:

```bash
make apply
```

To test the function, run the following command:

```bash
make test
```

To destroy the infrastructure, run the following command:

```bash
make destroy
```