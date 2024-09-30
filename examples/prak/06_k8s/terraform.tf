terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
    null = {
      source = "hashicorp/null"
    }
  }
  required_version = ">= 1.9"
}

data "yandex_client_config" "client" {}

provider "yandex" {
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = var.zone
}


provider "kubernetes" {
  host                   = module.k8s.external_v4_endpoint
  cluster_ca_certificate = module.k8s.cluster_ca_certificate
  token                  = data.yandex_client_config.client.iam_token
}
