data "yandex_compute_image" "ubuntu" {
  family = "ubuntu-2204-lts"
}


data "yandex_vpc_subnet" "default" {
  folder_id = var.folder_id
  name      = "default-ru-central1-a"
}

resource "yandex_compute_instance" "vm" {
  count       = 2
  name        = "test-vm${count.index}"
  platform_id = "standard-v3"
  zone        = "ru-central1-a"

  resources {
    cores  = 2
    memory = 4
  }

  boot_disk {
    initialize_params {
      image_id = data.yandex_compute_image.ubuntu.id
      size     = 20
    }
  }

  network_interface {
    index     = 1
    subnet_id = data.yandex_vpc_subnet.default.id
    nat = true
  }

  metadata = {
    user-data = <<-EOF
      #!/bin/bash
      apt-get update
      apt-get install -y nginx
      service nginx start
      sed -i -- "s/nginx/the first test/" /var/www/html/index.nginx-debian.html
    EOF
  }
}

resource "yandex_lb_target_group" "test" {
  name      = "test1"
  region_id = "ru-central1"

  dynamic "target" {
    for_each = yandex_compute_instance.vm
    content {
        subnet_id = data.yandex_vpc_subnet.default.id
        address   = target.value.network_interface.0.ip_address
    }
  }
}

resource "yandex_lb_network_load_balancer" "test" {
  name      = "test1-balancer"
  region_id = "ru-central1"
  listener {
    name = "my-listener"
    port = 80
    external_address_spec {
      ip_version = "ipv4"
    }
  }

  attached_target_group {
    target_group_id = yandex_lb_target_group.test.id

    healthcheck {
      name = "http"
      http_options {
        port = 80
        path = "/"
      }
    }
  }
}
