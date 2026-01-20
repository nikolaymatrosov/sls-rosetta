output "bucket_name" {
  description = "Name of the storage bucket"
  value       = yandex_storage_bucket.for_uploads.bucket
}

output "function_id" {
  description = "ID of the storage handler function"
  value       = yandex_function.storage_handler.id
}

output "trigger_id" {
  description = "ID of the storage trigger"
  value       = yandex_function_trigger.storage_trigger.id
}

output "access_key" {
  description = "S3 access key for manual testing"
  value       = yandex_iam_service_account_static_access_key.sa_storage_editor.access_key
  sensitive   = true
}

output "secret_key" {
  description = "S3 secret key for manual testing"
  value       = yandex_iam_service_account_static_access_key.sa_storage_editor.secret_key
  sensitive   = true
}

output "test_upload_command" {
  description = "AWS CLI command to upload test image"
  value       = <<-EOT
    # Upload a test image:
    aws s3 cp --endpoint-url=https://storage.yandexcloud.net \
      ./test-image.png \
      s3://${yandex_storage_bucket.for_uploads.bucket}/uploads/test-image.png
  EOT
}

output "test_list_command" {
  description = "AWS CLI command to list thumbnails"
  value       = <<-EOT
    # Check for thumbnail (wait a few seconds):
    aws s3 ls --endpoint-url=https://storage.yandexcloud.net \
      s3://${yandex_storage_bucket.for_uploads.bucket}/thumbnails/
  EOT
}
