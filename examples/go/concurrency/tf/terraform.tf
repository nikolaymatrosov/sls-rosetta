terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
      version = ">=0.123.0"
    }
  }
  required_version = ">= 1.8"

  # Настроить работу с s3
  backend "s3" {
    endpoints = {
      s3                = "https://storage.yandexcloud.net"
    }
    bucket = "cloud-terraform"
    region = "ru-central1"
    key    = "concurrency.tfstate"
#    /* shared_credentials_files = "storage.key" */
#    access_key     = "xxx"
#    secret_key     = "xxx"
    /* dynamodb_table = "ydb-tf-state-pr-vostok-1" */

    skip_region_validation      = true
    skip_credentials_validation = true
    skip_requesting_account_id  = true # Необходимая опция Terraform для версии 1.6.1 и старше.
    skip_s3_checksum            = true # Необходимая опция при описании бэкенда для Terraform версии 1.6.3 и старше.

  }
}

provider "yandex" {
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = var.zone
}