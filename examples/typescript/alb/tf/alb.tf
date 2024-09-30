resource "yandex_alb_backend_group" "s3-backend-group" {
  name      = "s3-group"

  http_backend {
    name = "test-http-backend"
    weight = 1
    storage_bucket = yandex_storage_bucket.test.bucket

    http2 = "true"
  }
}

resource "yandex_alb_http_router" "s3-router" {
  name   = "s3-http-router"
}



resource "yandex_alb_virtual_host" "ws-virtual-host" {
  name           = "ws-virtual-host"
  http_router_id = yandex_alb_http_router.s3-router.id
  route {
    name = "objstore-route"
    http_route {
      http_route_action {
        backend_group_id = yandex_alb_backend_group.s3-backend-group.id
        timeout          = "3s"
      }
    }
  }
}

resource "yandex_alb_load_balancer" "s3-balancer" {
  name = "s3-alb"

  network_id = yandex_vpc_network.s3.id

  allocation_policy {
    location {
      zone_id   = "ru-central1-a"
      subnet_id = yandex_vpc_subnet.s3-subnet-a.id
    }
  }

  listener {
    name = "s3-listener"
    endpoint {
      address {
        external_ipv4_address {
        }
      }
      ports = [8080]
    }
    http {
      handler {
        http_router_id = yandex_alb_http_router.s3-router.id
      }
    }
  }
}