# Service account for producer function
resource "yandex_iam_service_account" "producer_sa" {
  name      = "yds-producer-sa"
  folder_id = var.folder_id
}

# Service account for consumer function
resource "yandex_iam_service_account" "consumer_sa" {
  name      = "yds-consumer-sa"
  folder_id = var.folder_id
}

# Service account for trigger
resource "yandex_iam_service_account" "trigger_sa" {
  name      = "yds-trigger-sa"
  folder_id = var.folder_id
}

# IAM bindings for producer service account
resource "yandex_resourcemanager_folder_iam_binding" "producer_sa" {
  for_each = toset([
    "ydb.editor",
    "functions.functionInvoker",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.producer_sa.id}",
  ]
  sleep_after = 5
}

# IAM bindings for consumer service account
resource "yandex_resourcemanager_folder_iam_binding" "consumer_sa" {
  for_each = toset([
    "ydb.viewer",
    "functions.functionInvoker",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.consumer_sa.id}",
  ]
  sleep_after = 5
}

# IAM bindings for trigger service account
resource "yandex_resourcemanager_folder_iam_binding" "trigger_sa" {
  for_each = toset([
    "ydb.admin",
    "functions.functionInvoker",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.trigger_sa.id}",
  ]
  sleep_after = 5
} 