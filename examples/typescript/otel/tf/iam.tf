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
  description = "Service account for OpenTelemetry example functions (reads Lockbox)"
  folder_id   = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_member" "lockbox_payload_viewer" {
  folder_id = var.folder_id
  role      = "lockbox.payloadViewer"
  member    = "serviceAccount:${yandex_iam_service_account.otel_function_sa.id}"

  sleep_after = 5
}
