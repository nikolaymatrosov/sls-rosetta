resource "yandex_iam_service_account" "logging_transfer" {
  name      = "logging-transfer"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "logging_transfer" {
  for_each = toset([
    "yds.editor",
    "storage.editor",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.logging_transfer.id}",
  ]
  sleep_after = 5
}