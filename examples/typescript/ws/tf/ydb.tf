resource "yandex_ydb_database_serverless" "ws_db" {
  name        = "websocket-connections-db"
  folder_id   = var.folder_id
  location_id = "ru-central1"
}

data "dirhash_sha256" "migrations" {
  directory = "../migrations"
}

resource "null_resource" "run_migrations" {
  provisioner "local-exec" {
    command = <<-EOT
      export YDB_CONNECTION_STRING="${yandex_ydb_database_serverless.ws_db.ydb_full_endpoint}"
      goose -dir ../migrations ydb "$YDB_CONNECTION_STRING&token=$YC_TOKEN&go_query_mode=scripting&go_fake_tx=scripting&go_query_bind=declare,numeric" up
    EOT
  }

  triggers = {
    migrations_hash = data.dirhash_sha256.migrations.checksum
  }

  depends_on = [
    yandex_ydb_database_serverless.ws_db,
  ]
}

