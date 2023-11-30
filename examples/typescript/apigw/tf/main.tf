resource "null_resource" "build_typescript" {
  provisioner "local-exec" {
    command = "cd ../function && npm run build"
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
}

resource "yandex_function" "test_function" {
  name               = "api-gateway-demo"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs18"
  entrypoint         = "main.handler"
  memory             = "128"
  execution_timeout  = "10"
  service_account_id = yandex_iam_service_account.sa_serverless.id
  content {
    zip_filename = data.archive_file.function_files.output_path
  }
}



