resource "yandex_iam_service_account" "function_sa" {
  name      = "ydb-function-sa"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "function_sa" {
  for_each = toset([
    "ydb.viewer",
    "ydb.editor",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.function_sa.id}",
  ]
  sleep_after = 5
} 