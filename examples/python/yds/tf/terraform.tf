terraform {
  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = ">= 0.120"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.4"
    }
  }
  required_version = ">= 1.0"

  backend "local" {
    path = "../environment/terraform.tfstate"
  }
}

provider "yandex" {
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = var.zone
}
