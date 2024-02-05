resource "yandex_iam_service_account" "db_admin" {
  name      = "db-admin"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "db_admin" {
  for_each = toset([
    "ydb.admin",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.db_admin.id}",
  ]
  sleep_after = 5
}

resource "yandex_iam_service_account_static_access_key" "db_admin" {
  service_account_id = yandex_iam_service_account.db_admin.id
}

resource "yandex_iam_service_account" "lockbox_reader" {
  name      = "lockbox-reader"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "lockbox_reader" {
  for_each = toset([
    "lockbox.payloadViewer",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.lockbox_reader.id}",
  ]
  sleep_after = 5
}

