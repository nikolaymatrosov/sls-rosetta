resource "yandex_iam_service_account" "postbox_sender" {
  name      = "postbox-sender"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "sa_serverless" {
  for_each = toset([
    "postbox.sender",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.postbox_sender.id}",
  ]
  sleep_after = 5
}

resource "yandex_iam_service_account_static_access_key" "postbox_sender_key" {
  service_account_id = yandex_iam_service_account.postbox_sender.id
}