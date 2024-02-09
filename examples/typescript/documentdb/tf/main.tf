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
resource "yandex_function" "document-db-function" {
  name              = "document-db"
  user_hash         = data.archive_file.function_files.output_sha256
  runtime           = "nodejs18"
  entrypoint        = "main.handler"
  memory            = "128"
  execution_timeout = "10"
  content {
    zip_filename = data.archive_file.function_files.output_path
  }
  secrets {
    id                   = yandex_lockbox_secret.db-keys.id
    version_id           = yandex_lockbox_secret_version.db-keys.id
    key                  = "AWS_ACCESS_KEY_ID"
    environment_variable = "AWS_ACCESS_KEY_ID"
  }
  secrets {
    id                   = yandex_lockbox_secret.db-keys.id
    version_id           = yandex_lockbox_secret_version.db-keys.id
    key                  = "AWS_SECRET_ACCESS_KEY"
    environment_variable = "AWS_SECRET_ACCESS_KEY"
  }
  environment = {
    ENDPOINT = yandex_ydb_database_serverless.db.document_api_endpoint
  }
  service_account_id = yandex_iam_service_account.lockbox_reader.id

  depends_on = [
    yandex_lockbox_secret.db-keys,
    yandex_lockbox_secret_version.db-keys,
    yandex_ydb_database_serverless.db,
    yandex_iam_service_account.lockbox_reader
  ]
}

// IAM binding for making function public
resource "yandex_function_iam_binding" "test_function_binding" {
  function_id = yandex_function.document-db-function.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}



