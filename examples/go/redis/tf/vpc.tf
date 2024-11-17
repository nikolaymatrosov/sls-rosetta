data "yandex_vpc_network" "default" {
    folder_id = var.folder_id
    name      = "default"
}


locals {
  zones = {
    "ru-central1-a" = "10.0.0.0/24",
    "ru-central1-b" = "10.1.0.0/24",
    "ru-central1-d" = "10.2.0.0/24",
  }
}

resource "yandex_vpc_subnet" "redis" {
  for_each = local.zones
  zone       = each.key
  network_id = data.yandex_vpc_network.default.id
  v4_cidr_blocks = [each.value]
}