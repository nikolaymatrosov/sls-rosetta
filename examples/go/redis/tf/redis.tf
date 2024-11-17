
resource "yandex_lockbox_secret" "redis" {
  name = "redis-password"
}

resource "random_password" "redis" {
  length = 24
  override_special = "@=+?*.,!&#$^<>_-"
}

resource "yandex_lockbox_secret_version_hashed" "password_version" {
  secret_id = yandex_lockbox_secret.redis.id
  key_1 = "password"
  text_value_1 = random_password.redis.result
}

resource "yandex_mdb_redis_cluster" "demo" {
  name        = "test"
  environment = "PRESTABLE"
  network_id  = data.yandex_vpc_network.default.id

  config {
    password = random_password.redis.result
    version  = "7.2"

  }

  access {
    web_sql = true
  }

  resources {
    resource_preset_id = "hm2.nano"
    disk_size          = 16
  }

  dynamic "host" {
    for_each = local.zones
    content {
      zone      = host.key
      subnet_id = yandex_vpc_subnet.redis[host.key].id
    }
  }

  maintenance_window {
    type = "ANYTIME"
  }
}
