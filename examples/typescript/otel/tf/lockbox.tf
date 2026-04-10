resource "yandex_lockbox_secret" "monium_api_key" {
  name      = "monium-api-key"
  folder_id = var.folder_id
}

resource "yandex_lockbox_secret_version" "monium_api_key" {
  secret_id = yandex_lockbox_secret.monium_api_key.id
  entries {
    key        = "MONIUM_API_KEY"
    text_value = yandex_iam_service_account_api_key.otel_monium_api_key.secret_key
  }
}
