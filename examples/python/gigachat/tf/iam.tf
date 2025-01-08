resource "yandex_iam_service_account" "lockbox" {
  name      = "lockbox-reader"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "lockbox" {
  for_each = toset([
    "lockbox.payloadViewer",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.lockbox.id}",
  ]
  sleep_after = 5
}