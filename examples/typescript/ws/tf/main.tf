resource "null_resource" "build_typescript" {
  provisioner "local-exec" {
    command = "cd ../function && npm install && npm run build"
  }
  triggers = {
    always_run = timestamp()
  }
}

data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../dist"
  type        = "zip"
  depends_on  = [
    null_resource.build_typescript
  ]
  excludes = [ "node_modules" ]
}


resource "yandex_function" "ws_handler" {
  name               = "websocket-broadcast-handler"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs22"
  entrypoint         = "server/main.handler"
  memory             = "256"
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.ws_function_sa.id

  content {
    zip_filename = data.archive_file.function_files.output_path
  }

  environment = {
    YDB_CONNECTION_STRING = yandex_ydb_database_serverless.ws_db.ydb_full_endpoint
  }

  depends_on = [
    yandex_ydb_database_serverless.ws_db,
    yandex_iam_service_account.ws_function_sa,
    yandex_resourcemanager_folder_iam_binding.ws_function_sa,
    null_resource.run_migrations,
  ]
}

resource "yandex_function" "ws_handler_v2" {
  name               = "websocket-broadcast-handler-v2"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs22"
  entrypoint         = "server/main_v2.handler"
  memory             = "256"
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.ws_function_sa.id

  content {
    zip_filename = data.archive_file.function_files.output_path
  }

  environment = {
    YDB_CONNECTION_STRING = yandex_ydb_database_serverless.ws_db.ydb_full_endpoint
    BROADCAST_TOPIC       = yandex_ydb_topic.broadcast_topic.name
  }

  depends_on = [
    yandex_ydb_database_serverless.ws_db,
    yandex_ydb_topic.broadcast_topic,
    yandex_iam_service_account.ws_function_sa,
    yandex_resourcemanager_folder_iam_binding.ws_function_sa,
    null_resource.run_migrations,
  ]
}