resource "yandex_iam_service_account" "sa_ymq_creator" {
  name      = "sa-ymq-creator"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "sa_ymq_creator" {
  for_each = toset([
    "ymq.admin",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.sa_ymq_creator.id}",
  ]
  sleep_after = 5
}

resource "yandex_iam_service_account_static_access_key" "sa_ymq_creator" {
  service_account_id = yandex_iam_service_account.sa_ymq_creator.id
}

resource "yandex_iam_service_account" "trigger_sa" {
  name      = "ymq-trigger-sa"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "trigger_sa" {
  for_each = toset([
    "ymq.reader",
    "ymq.writer",
    "functions.functionInvoker"
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.trigger_sa.id}",
  ]
}

resource "yandex_iam_service_account" "ymq_writer" {
  name      = "ymq-writer"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "ymq_writer" {
  for_each = toset([
    "ymq.writer",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.ymq_writer.id}",
  ]

}

resource "yandex_iam_service_account_static_access_key" "ymq_writer" {
  service_account_id = yandex_iam_service_account.ymq_writer.id
}

// Currently, ymq_reader is used only for tests
// You can remove it and its dependencies if you don't need it
resource "yandex_iam_service_account" "ymq_reader" {
  name      = "ymq-reader"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "ymq_reader" {
  for_each = toset([
    "ymq.reader",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.ymq_reader.id}",
  ]
}

resource "yandex_iam_service_account_static_access_key" "ymq_reader" {
  service_account_id = yandex_iam_service_account.ymq_reader.id
}


