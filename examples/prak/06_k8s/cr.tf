resource "yandex_container_registry" "praktikum" {
  name      = "praktikum"
  folder_id = var.folder_id
}

resource "yandex_container_repository" "praktikum" {
  name = "${yandex_container_registry.praktikum.id}/nginx"
}

resource "yandex_iam_service_account" "praktikum" {
  name        = "praktikum"
  folder_id   = var.folder_id
  description = "Service account for praktikum"
}

resource "yandex_iam_service_account_key" "praktikum" {
  service_account_id = yandex_iam_service_account.praktikum.id
}

resource "yandex_container_registry_iam_binding" "puller" {
  registry_id = yandex_container_registry.praktikum.id
  role        = "container-registry.images.pusher"

  members = [
    "serviceAccount:${yandex_iam_service_account.praktikum.id}",
  ]
}

locals {
  json_key = jsonencode({
    "id" : yandex_iam_service_account_key.praktikum.id
    "service_account_id" : yandex_iam_service_account_key.praktikum.service_account_id
    "created_at" : yandex_iam_service_account_key.praktikum.created_at
    "key_algorithm" : yandex_iam_service_account_key.praktikum.key_algorithm
    "public_key" : yandex_iam_service_account_key.praktikum.public_key
    "private_key" : yandex_iam_service_account_key.praktikum.private_key
  })
  image = "cr.yandex/${yandex_container_repository.praktikum.name}:latest"
}

resource "null_resource" "nginx" {
  provisioner "local-exec" {
    command = <<-EOF
        docker pull nginx:latest --platform amd64
        docker login -u json_key -p ${local.json_key} cr.yandex
        docker tag nginx:latest ${local.image}
        docker push ${local.image}
    EOF
  }
}