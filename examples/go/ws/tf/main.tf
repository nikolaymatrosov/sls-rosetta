locals {
  function_dir = "${path.module}/../function/server"
}

# Archive the function code
data "archive_file" "function_files" {
  type        = "zip"
  output_path = "${path.module}/function.zip"
  source_dir  = local.function_dir
  excludes    = [
    "go.sum",
    ".gitignore",
  ]
}

# Create the WebSocket handler function
resource "yandex_function" "ws_handler" {
  name               = var.function_name
  runtime            = "golang123"
  entrypoint         = "index.WebSocketEventHandler"
  memory             = "256"
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.ws_sa.id
  user_hash          = data.archive_file.function_files.output_sha256

  environment = {
    YDB_CONNECTION_STRING = yandex_ydb_database_serverless.ws_database.ydb_full_endpoint
    YDB_DATABASE          = yandex_ydb_database_serverless.ws_database.database_path
    BROADCAST_TOPIC       = "${yandex_ydb_database_serverless.ws_database.database_path}/${var.topic_name}"
  }

  content {
    zip_filename = data.archive_file.function_files.output_path
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.sa_ydb_editor,
    yandex_resourcemanager_folder_iam_member.sa_websocket_writer,
    yandex_resourcemanager_folder_iam_member.sa_websocket_broadcaster,
  ]
}

# Grant public access to the function
resource "yandex_function_iam_binding" "function_invoker" {
  function_id = yandex_function.ws_handler.id
  role        = "serverless.functions.invoker"

  members = [
    "system:allUsers",
  ]
}
