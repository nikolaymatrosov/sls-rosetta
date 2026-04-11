resource "yandex_iam_service_account" "otel_monium_sa" {
  name        = "otel-monium-sa"
  description = "Service account with Monium traces writer role"
  folder_id   = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_member" "otel_traces_writer" {
  folder_id = var.folder_id
  role      = "monium.traces.writer"
  member    = "serviceAccount:${yandex_iam_service_account.otel_monium_sa.id}"

  sleep_after = 5
}

resource "yandex_iam_service_account_api_key" "otel_monium_api_key" {
  service_account_id = yandex_iam_service_account.otel_monium_sa.id
  description        = "API key for Monium traces ingestion"
  scopes = [
    "yc.monium.traces.write",
  ]
  lifecycle {
    ignore_changes = [ scope ]
  }
}

resource "yandex_iam_service_account" "otel_function_sa" {
  name        = "otel-function-sa"
  description = "Service account for OpenTelemetry example functions (reads Lockbox, writes MQ)"
  folder_id   = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_member" "lockbox_payload_viewer" {
  folder_id = var.folder_id
  role      = "lockbox.payloadViewer"
  member    = "serviceAccount:${yandex_iam_service_account.otel_function_sa.id}"

  sleep_after = 5
}

resource "yandex_resourcemanager_folder_iam_member" "otel_function_sa_ymq_writer" {
  folder_id = var.folder_id
  role      = "ymq.writer"
  member    = "serviceAccount:${yandex_iam_service_account.otel_function_sa.id}"

  sleep_after = 5
}

resource "yandex_iam_service_account_static_access_key" "otel_function_sa" {
  service_account_id = yandex_iam_service_account.otel_function_sa.id
  description        = "Static key for YMQ operations"
}

# Dedicated SA for queue provisioning (needs ymq.admin to create queues)
resource "yandex_iam_service_account" "otel_mq_admin_sa" {
  name        = "otel-mq-admin-sa"
  description = "Service account for provisioning YMQ queues"
  folder_id   = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_member" "otel_mq_admin_sa" {
  folder_id = var.folder_id
  role      = "ymq.admin"
  member    = "serviceAccount:${yandex_iam_service_account.otel_mq_admin_sa.id}"

  sleep_after = 5
}

resource "yandex_iam_service_account_static_access_key" "otel_mq_admin_sa" {
  service_account_id = yandex_iam_service_account.otel_mq_admin_sa.id
  description        = "Static key for queue provisioning"

  depends_on = [yandex_resourcemanager_folder_iam_member.otel_mq_admin_sa]
}

# Trigger SA — reads from YMQ and invokes the inventory function
resource "yandex_iam_service_account" "otel_mq_trigger_sa" {
  name        = "otel-mq-trigger-sa"
  description = "Service account for YMQ trigger invoking inventory function"
  folder_id   = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_member" "otel_mq_trigger_sa_roles" {
  for_each = toset([
    "ymq.reader",
    "functions.functionInvoker",
  ])
  
  folder_id = var.folder_id
  role      = each.value
  member    = "serviceAccount:${yandex_iam_service_account.otel_mq_trigger_sa.id}"

  sleep_after = 5
}

