resource "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../function"
  type        = "zip"
}

resource "yandex_function" "ydb_function" {
  name               = "ydb-demo"
  user_hash          = archive_file.function_files.output_sha256
  runtime            = "golang123"
  entrypoint         = "main.Handler"
  memory             = "128"
  execution_timeout  = "10"
  service_account_id = yandex_iam_service_account.function_sa.id
  
  content {
    zip_filename = archive_file.function_files.output_path
  }
  
  environment = {
    YDB_DATABASE = yandex_ydb_database_serverless.db.database_path
    YDB_ENDPOINT = yandex_ydb_database_serverless.db.ydb_api_endpoint
  }
  
  depends_on = [
    yandex_ydb_database_serverless.db,
    yandex_iam_service_account.function_sa,
    yandex_resourcemanager_folder_iam_binding.function_sa,
  ]
}

// IAM binding for making function public
resource "yandex_function_iam_binding" "function_binding" {
  function_id = yandex_function.ydb_function.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
} 