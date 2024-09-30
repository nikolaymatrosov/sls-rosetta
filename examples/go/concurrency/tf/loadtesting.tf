data "yandex_vpc_subnet" "subnet-a" {
  name = "default-ru-central1-a"
}

resource "yandex_iam_service_account" "load_testing" {
  name      = "load-testing"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "load_testing" {
  for_each = toset([
    "loadtesting.generatorClient",
    "logging.writer",
    "storage.editor",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.load_testing.id}",
  ]
  sleep_after = 5
}

resource "yandex_loadtesting_agent" "my-agent" {
  name        = "concurrency-tester"
  description = "2 core 4 GB RAM agent"
  folder_id   = var.folder_id


  compute_instance {
    zone_id            = var.zone
    service_account_id = yandex_iam_service_account.load_testing.id
    resources {
      memory = 4
      cores  = 2
    }
    boot_disk {
      initialize_params {
        size = 15
      }
      auto_delete = true
    }
    network_interface {
      subnet_id = data.yandex_vpc_subnet.subnet-a.id
    }
  }
}

resource "yandex_logging_group" "load_testing" {
  name      = "load-testing-logs"
  folder_id = var.folder_id
}