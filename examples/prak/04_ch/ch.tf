locals {
  user     = "user"
  password = "your_password"
  db_name  = "db_name"
}

resource "yandex_mdb_clickhouse_cluster" "praktikum" {
  name        = "praktikum"
  environment = "PRESTABLE"
  network_id  = data.yandex_vpc_subnet.default.network_id

  clickhouse {
    resources {
      resource_preset_id = "b3-c1-m4"
      disk_type_id       = "network-ssd"
      disk_size          = 10
    }
  }

  database {
    name = local.db_name
  }

  user {
    name     = "user"
    password = "your_password"
    permission {
      database_name = local.db_name
    }

    settings {
      max_memory_usage_for_user               = 1000000000
      read_overflow_mode                      = "throw"
      output_format_json_quote_64bit_integers = true
    }
  }

  host {
    type             = "CLICKHOUSE"
    zone             = var.zone
    subnet_id        = data.yandex_vpc_subnet.default.id
    assign_public_ip = true
  }

}

data "yandex_vpc_subnet" "default" {
  folder_id = var.folder_id
  name      = "default-ru-central1-a"
}


output "clickhouse" {
  //clickhouse-client --host=<your_host> --port=<your_port> --user=<your_username> --password=<your_password>
  value = "clickhouse client --host=${yandex_mdb_clickhouse_cluster.praktikum.host[0].fqdn} --port=8443 --user=${local.user} --password=${local.password}"
}

output "cluster_endpoint" {
  value = yandex_mdb_clickhouse_cluster.praktikum.host[0].fqdn
}

output "db_name" {
  value = local.db_name
}

output "user" {
  value = local.user
}

output "db_password" {
  value = local.password
}