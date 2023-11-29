output "bucket" {
  value = yandex_storage_bucket.for-uploads.bucket
}
output "bucket_for_function" {
  value = yandex_storage_bucket.for-deploy.bucket
}
output "sa_storage_editor_access_key" {
  value     = yandex_iam_service_account_static_access_key.sa_storage_editor.access_key
  sensitive = true
}

output "sa_storage_editor_secret_key" {
  value     = yandex_iam_service_account_static_access_key.sa_storage_editor.secret_key
  sensitive = true
}
