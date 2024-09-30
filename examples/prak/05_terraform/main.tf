resource "yandex_vpc_network" "test" {
  name      = "from-terraform-network"
  folder_id = var.folder_id
}

resource "yandex_vpc_subnet" "test" {
  name           = "from-terraform-subnet"
  network_id     = yandex_vpc_network.test.id
  folder_id      = var.folder_id
  zone           = "ru-central1-a"
  v4_cidr_blocks = ["10.2.0.0/16"]
}

data "yandex_compute_image" "ubuntu" {
  family = "ubuntu-2204-lts"
}

resource "yandex_compute_instance" "test" {
  name        = "from-terraform-vm"
  folder_id   = var.folder_id
  zone        = "ru-central1-a"
  platform_id = "standard-v1"
  resources {
    cores  = 2
    memory = 2
  }

  boot_disk {
    initialize_params {
      image_id = data.yandex_compute_image.ubuntu.id
    }
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.test.id
    nat = true
  }
}

#resource "yandex_mdb_postgresql_cluster" "test" {
#  folder_id   = var.folder_id
#  environment = "PRODUCTION"
#  name        = "test-vm"
#  network_id  = yandex_vpc_network.test.id
#  config {
#    version = "16"
#    resources {
#      disk_size          = 10
#      disk_type_id       = "network-ssd"
#      resource_preset_id = "s1.micro"
#    }
#  }
#  database {
#    name  = "test-db"
#    owner = "user"
#  }
#  host {
#    zone = "ru-central1-a"
#  }
#  user {
#    name     = "user"
#    password = "password"
#  }
#}
