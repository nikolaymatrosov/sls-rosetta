terraform {
  backend "s3" {
    region = "ru-central1"
    bucket = "nikolaymatrosov-tfstate"
    key    = "terraform.tfstate"

    dynamodb_table = "table732"

    endpoint          = "https://storage.yandexcloud.net"
    dynamodb_endpoint = "https://docapi.serverless.yandexcloud.net/ru-central1/b1glihojf6il6g7lvk98/etno3lual1aqavur902l"

    skip_credentials_validation = true
    skip_region_validation      = true
    #    skip_requesting_account_id  = true
    #    skip_s3_checksum            = true
  }

  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
  required_version = "1.5.7"
}

provider "yandex" {
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = var.zone
}