resource "yandex_iam_service_account" "sa_storage_editor" {
  name      = "sa-storage-admin"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "sa_storage_editor" {
  for_each = toset([
    "storage.editor",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.sa_storage_editor.id}",
  ]
  sleep_after = 5
}

resource "yandex_iam_service_account_static_access_key" "sa_storage_editor" {
  service_account_id = yandex_iam_service_account.sa_storage_editor.id
}

resource "yandex_iam_service_account" "trigger_sa" {
  name      = "ymq-trigger-sa"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "trigger_sa" {
  for_each = toset([
    "functions.functionInvoker"
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.trigger_sa.id}",
  ]
}
