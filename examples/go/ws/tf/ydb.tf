# Create YDB serverless database
resource "yandex_ydb_database_serverless" "ws_database" {
  name      = var.database_name
  folder_id = var.folder_id

  serverless_database {
    enable_throttling_rcu_limit = false
    provisioned_rcu_limit       = 10
    storage_size_limit          = 50
    throttling_rcu_limit        = 0
  }
}

data "dirhash_sha256" "migrations" {
  directory = "../migrations"
}

resource "null_resource" "run_migrations" {
  provisioner "local-exec" {
    command = <<-EOT
      export YDB_CONNECTION_STRING="${yandex_ydb_database_serverless.ws_database.ydb_full_endpoint}"
      export IAM_TOKEN="$(yc iam create-token)"
      goose -dir ../migrations ydb "$YDB_CONNECTION_STRING&token=$IAM_TOKEN&go_query_mode=scripting&go_fake_tx=scripting&go_query_bind=declare,numeric" up
    EOT
  }

  triggers = {
    migrations_hash = data.dirhash_sha256.migrations.checksum
  }

  depends_on = [
    yandex_ydb_database_serverless.ws_database,
  ]
}

# Create YDB topic for broadcasting
resource "yandex_ydb_topic" "broadcast_topic" {
  database_endpoint = yandex_ydb_database_serverless.ws_database.ydb_full_endpoint
  name              = var.topic_name

  supported_codecs = ["raw"]

  depends_on = [
    yandex_ydb_database_serverless.ws_database,
  ]
}

# Data Streams trigger for broadcasting WebSocket messages
# Note: Using null_resource because Terraform provider doesn't support gateway_websocket_broadcast yet
resource "null_resource" "broadcast_trigger" {
  provisioner "local-exec" {
    command = <<-EOT
      yc serverless trigger create yds ws-go-broadcast-trigger \
        --folder-id ${var.folder_id} \
        --description "Triggers on messages in broadcast topic to send them to WebSocket connections" \
        --database ${yandex_ydb_database_serverless.ws_database.database_path} \
        --stream ${yandex_ydb_topic.broadcast_topic.name} \
        --stream-service-account-id ${yandex_iam_service_account.ws_sa.id} \
        --batch-size 1b \
        --batch-cutoff 1s \
        --gateway-id ${yandex_api_gateway.ws_gateway.id} \
        --gateway-websocket-broadcast-path /ws \
        --gateway-websocket-broadcast-service-account-id ${yandex_iam_service_account.ws_sa.id}
    EOT
  }

  provisioner "local-exec" {
    when    = destroy
    command = "yc serverless trigger delete ws-go-broadcast-trigger --folder-id ${self.triggers.folder_id} || true"
  }

  triggers = {
    folder_id          = var.folder_id
    database_path      = yandex_ydb_database_serverless.ws_database.database_path
    topic_name         = yandex_ydb_topic.broadcast_topic.name
    gateway_id         = yandex_api_gateway.ws_gateway.id
    service_account_id = yandex_iam_service_account.ws_sa.id
  }

  depends_on = [
    yandex_ydb_topic.broadcast_topic,
    yandex_api_gateway.ws_gateway,
    yandex_resourcemanager_folder_iam_member.sa_yds_admin,
    yandex_resourcemanager_folder_iam_member.sa_websocket_broadcaster,
  ]
}
