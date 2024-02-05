data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../ApiGW/ApiGWHandler"
  type        = "zip"
}

resource "yandex_function" "test_function" {
  name               = "api-gateway-demo"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "dotnet8"
  entrypoint         = "ApiGWHandler.Handler"
  memory             = "128"
  execution_timeout  = "10"
  service_account_id = yandex_iam_service_account.sa_serverless.id
  content {
    zip_filename = data.archive_file.function_files.output_path
  }
}
