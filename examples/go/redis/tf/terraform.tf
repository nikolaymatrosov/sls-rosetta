terraform {
  backend "local" {
    path = "../environment/terraform.tfstate"
  }
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
      version = ">= 0.104"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.6.2"
    }
  }
  required_version = ">= 1.0"
}

provider "yandex" {
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = var.zone
}