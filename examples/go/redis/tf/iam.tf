resource "yandex_iam_service_account" "handler" {
  name      = "lockbox-sa"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "handler" {
  for_each = toset([
    "lockbox.payloadViewer",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.handler.id}",
  ]
  sleep_after = 5
}
