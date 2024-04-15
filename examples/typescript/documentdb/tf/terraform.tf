terraform {
  backend "local" {
    path = "../environment/terraform.tfstate"
  }
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
    null = {
      source = "hashicorp/null"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  required_version = ">= 0.13"
}

provider "yandex" {
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = var.zone
}

provider "aws" {
  region = "ru-central1"
  endpoints {
    dynamodb = yandex_ydb_database_serverless.db.document_api_endpoint
  }
  access_key = yandex_iam_service_account_static_access_key.db_admin.access_key
  secret_key = yandex_iam_service_account_static_access_key.db_admin.secret_key
  skip_credentials_validation = true
  skip_metadata_api_check = true
  skip_region_validation = true
  skip_requesting_account_id = true
}