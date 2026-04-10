resource "null_resource" "build_typescript" {
  provisioner "local-exec" {
    command = "cd ../function && npm ci && npm run build"
  }
  triggers = {
    always_run = timestamp()
  }
}

data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../dist"
  type        = "zip"
  depends_on = [
    null_resource.build_typescript
  ]
}

resource "yandex_function" "otel_callee" {
  name               = "otel-callee-ts"
  description        = "OpenTelemetry example: callee function"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs22"
  entrypoint         = "callee.handler"
  memory             = 256
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.otel_function_sa.id

  environment = {
    OTEL_SERVICE_NAME = "payment-service"
  }

  secrets {
    id                   = yandex_lockbox_secret.monium_api_key.id
    version_id           = yandex_lockbox_secret_version.monium_api_key.id
    key                  = "MONIUM_API_KEY"
    environment_variable = "MONIUM_API_KEY"
  }

  content {
    zip_filename = data.archive_file.function_files.output_path
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.lockbox_payload_viewer,
    yandex_lockbox_secret_version.monium_api_key,
  ]
}

resource "yandex_function" "otel_caller" {
  name               = "otel-caller-ts"
  description        = "OpenTelemetry example: caller function"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs22"
  entrypoint         = "caller.handler"
  memory             = 256
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.otel_function_sa.id

  environment = {
    OTEL_SERVICE_NAME = "order-service"
    CALLEE_URL        = "https://functions.yandexcloud.net/${yandex_function.otel_callee.id}"
  }

  secrets {
    id                   = yandex_lockbox_secret.monium_api_key.id
    version_id           = yandex_lockbox_secret_version.monium_api_key.id
    key                  = "MONIUM_API_KEY"
    environment_variable = "MONIUM_API_KEY"
  }

  content {
    zip_filename = data.archive_file.function_files.output_path
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.lockbox_payload_viewer,
    yandex_lockbox_secret_version.monium_api_key,
  ]
}

resource "yandex_function_iam_binding" "caller_public" {
  function_id = yandex_function.otel_caller.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}

resource "yandex_function_iam_binding" "callee_public" {
  function_id = yandex_function.otel_callee.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}
