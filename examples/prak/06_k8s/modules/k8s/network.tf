resource "yandex_vpc_network" "default" {
  name = var.network_name
}

resource "yandex_vpc_subnet" "default" {
  for_each = {
    a = "10.128.0.0/24",
    b = "10.129.0.0/24",
    d = "10.131.0.0/24",
  }
  name       = "default-${each.key}"
  network_id = yandex_vpc_network.default.id
  zone       = "ru-central1-${each.key}"

  v4_cidr_blocks = [
    each.value
  ]
  route_table_id = yandex_vpc_route_table.default.id
}

resource "yandex_vpc_gateway" "default" {
  name = "default"
  shared_egress_gateway {}
}

resource "yandex_vpc_route_table" "default" {
  name       = "default"
  network_id = yandex_vpc_network.default.id
  static_route {
    destination_prefix = "0.0.0.0/0"
    gateway_id         = yandex_vpc_gateway.default.id
  }
}


