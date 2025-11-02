# Service Account for Producer Function
resource "yandex_iam_service_account" "producer_sa" {
  name        = "yds-producer-sa-ts"
  description = "Service account for YDS producer function (TypeScript)"
  folder_id   = var.folder_id
}

# Service Account for Consumer Function
resource "yandex_iam_service_account" "consumer_sa" {
  name        = "yds-consumer-sa-ts"
  description = "Service account for YDS consumer function (TypeScript)"
  folder_id   = var.folder_id
}

# Service Account for Trigger
resource "yandex_iam_service_account" "trigger_sa" {
  name        = "yds-trigger-sa-ts"
  description = "Service account for YDS trigger (TypeScript)"
  folder_id   = var.folder_id
}

# IAM Bindings for Producer
resource "yandex_resourcemanager_folder_iam_member" "producer_roles" {
  for_each = toset([
    "ydb.editor",
    "functions.functionInvoker"
  ])

  folder_id = var.folder_id
  role      = each.value
  member    = "serviceAccount:${yandex_iam_service_account.producer_sa.id}"

  sleep_after = 5
}

# IAM Bindings for Consumer
resource "yandex_resourcemanager_folder_iam_member" "consumer_roles" {
  for_each = toset([
    "ydb.viewer",
    "functions.functionInvoker"
  ])

  folder_id = var.folder_id
  role      = each.value
  member    = "serviceAccount:${yandex_iam_service_account.consumer_sa.id}"

  sleep_after = 5
}

# IAM Bindings for Trigger
resource "yandex_resourcemanager_folder_iam_member" "trigger_roles" {
  for_each = toset([
    "ydb.admin",
    "functions.functionInvoker"
  ])

  folder_id = var.folder_id
  role      = each.value
  member    = "serviceAccount:${yandex_iam_service_account.trigger_sa.id}"

  sleep_after = 5
}

# Make producer function publicly accessible
resource "yandex_function_iam_binding" "producer_public" {
  function_id = yandex_function.producer.id
  role        = "functions.functionInvoker"
  members = [
    "system:allUsers"
  ]
}
