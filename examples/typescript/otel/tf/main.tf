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

resource "yandex_function" "otel_payment" {
  name               = "otel-payment-ts"
  description        = "OpenTelemetry example: payment service"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs22"
  entrypoint         = "payment.handler"
  memory             = 256
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.otel_function_sa.id

  environment = {
    OTEL_SERVICE_NAME = "payment-service"
    FOLDER_ID         = var.folder_id
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

resource "yandex_function" "otel_orders" {
  name               = "otel-orders-ts"
  description        = "OpenTelemetry example: order service — calls payment and enqueues to YMQ"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs22"
  entrypoint         = "orders.handler"
  memory             = 256
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.otel_function_sa.id

  environment = {
    OTEL_SERVICE_NAME     = "order-service"
    FOLDER_ID             = var.folder_id
    CALLEE_URL            = "https://functions.yandexcloud.net/${yandex_function.otel_payment.id}"
    ORDERS_QUEUE_URL      = yandex_message_queue.orders_queue.id
    AWS_ACCESS_KEY_ID     = yandex_iam_service_account_static_access_key.otel_function_sa.access_key
    AWS_SECRET_ACCESS_KEY = yandex_iam_service_account_static_access_key.otel_function_sa.secret_key
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
    yandex_resourcemanager_folder_iam_member.otel_function_sa_ymq_writer,
    yandex_lockbox_secret_version.monium_api_key,
  ]
}

resource "yandex_function" "otel_inventory" {
  name               = "otel-inventory-ts"
  description        = "OpenTelemetry example: inventory service — consumes orders from YMQ"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs22"
  entrypoint         = "inventory.handler"
  memory             = 256
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.otel_function_sa.id

  environment = {
    OTEL_SERVICE_NAME = "inventory-service"
    FOLDER_ID         = var.folder_id
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

resource "yandex_function_trigger" "orders_trigger" {
  name        = "otel-orders-trigger"
  description = "Triggers inventory function on new orders in YMQ"

  message_queue {
    queue_id           = yandex_message_queue.orders_queue.arn
    service_account_id = yandex_iam_service_account.otel_mq_trigger_sa.id
    batch_size         = "5"
    batch_cutoff       = "5"
  }

  function {
    id                 = yandex_function.otel_inventory.id
    service_account_id = yandex_iam_service_account.otel_mq_trigger_sa.id
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.otel_mq_trigger_sa_roles
  ]
}

resource "yandex_function_iam_binding" "orders_public" {
  function_id = yandex_function.otel_orders.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}

resource "yandex_function_iam_binding" "payment_public" {
  function_id = yandex_function.otel_payment.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}
