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

# YDB Topic for broadcasting messages via Data Streams trigger
resource "yandex_ydb_topic" "broadcast_topic" {
  database_endpoint = yandex_ydb_database_serverless.ws_db.ydb_full_endpoint
  name              = "broadcast-topic"

  supported_codecs = ["raw"]

  consumer {
    name             = "broadcast-consumer"
    starting_message_timestamp_ms = 0
    supported_codecs = ["raw"]
  }

  depends_on = [
    yandex_ydb_database_serverless.ws_db,
  ]
}

# Data Streams trigger for broadcasting WebSocket messages
# Note: Using null_resource because Terraform provider doesn't support gateway_websocket_broadcast yet
resource "null_resource" "broadcast_trigger" {
  provisioner "local-exec" {
    command = <<-EOT
      yc serverless trigger create yds websocket-broadcast-trigger \
        --folder-id ${var.folder_id} \
        --description "Triggers on messages in broadcast topic to send them to WebSocket connections" \
        --database ${yandex_ydb_database_serverless.ws_db.database_path} \
        --stream ${yandex_ydb_topic.broadcast_topic.name} \
        --stream-service-account-id ${yandex_iam_service_account.ws_function_sa.id} \
        --batch-size 1b \
        --batch-cutoff 1s \
        --gateway-id ${yandex_api_gateway.ws_gateway.id} \
        --gateway-websocket-broadcast-path /ws_v2 \
        --gateway-websocket-broadcast-service-account-id ${yandex_iam_service_account.ws_function_sa.id}
    EOT
  }

  provisioner "local-exec" {
    when    = destroy
    command = "yc serverless trigger delete websocket-broadcast-trigger --folder-id ${self.triggers.folder_id} || true"
  }

  triggers = {
    folder_id          = var.folder_id
    database_path      = yandex_ydb_database_serverless.ws_db.database_path
    topic_name         = yandex_ydb_topic.broadcast_topic.name
    gateway_id         = yandex_api_gateway.ws_gateway.id
    service_account_id = yandex_iam_service_account.ws_function_sa.id
  }

  depends_on = [
    yandex_ydb_topic.broadcast_topic,
    yandex_api_gateway.ws_gateway,
    yandex_resourcemanager_folder_iam_binding.ws_function_sa,
  ]
}

