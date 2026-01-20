# Service Account for Storage Operations
resource "yandex_iam_service_account" "sa_storage_editor" {
  name        = "storage-handler-sa-python"
  description = "Service account for storage operations and function execution"
  folder_id   = var.folder_id
}

# IAM Roles for Storage Service Account
resource "yandex_resourcemanager_folder_iam_member" "sa_storage_editor" {
  for_each = toset([
    "storage.editor",
  ])

  folder_id = var.folder_id
  role      = each.value
  member    = "serviceAccount:${yandex_iam_service_account.sa_storage_editor.id}"

  sleep_after = 5
}

# Static Access Key for S3 operations
resource "yandex_iam_service_account_static_access_key" "sa_storage_editor" {
  service_account_id = yandex_iam_service_account.sa_storage_editor.id
  description        = "Static access key for S3 operations"
}

# Service Account for Trigger
resource "yandex_iam_service_account" "trigger_sa" {
  name        = "storage-trigger-sa-python"
  description = "Service account for storage trigger"
  folder_id   = var.folder_id
}

# IAM Roles for Trigger Service Account
resource "yandex_resourcemanager_folder_iam_member" "trigger_sa" {
  for_each = toset([
    "functions.functionInvoker"
  ])

  folder_id = var.folder_id
  role      = each.value
  member    = "serviceAccount:${yandex_iam_service_account.trigger_sa.id}"

  sleep_after = 5
}
