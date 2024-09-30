resource "yandex_vpc_network" "s3" {
  name = "s3"
}

resource "yandex_vpc_subnet" "s3-subnet-a" {
  name           = "s3-subnet-a"
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network.s3.id
  v4_cidr_blocks = ["10.240.1.0/24"]
}

resource "yandex_vpc_subnet" "s3-subnet-b" {
  name           = "s3-subnet-b"
  zone           = "ru-central1-b"
  network_id     = yandex_vpc_network.s3.id
  v4_cidr_blocks = ["10.240.2.0/24"]
}

resource "yandex_vpc_subnet" "s3-subnet-c" {
  name           = "s3-subnet-c"
  zone           = "ru-central1-c"
  network_id     = yandex_vpc_network.s3.id
  v4_cidr_blocks = ["10.240.3.0/24"]
}