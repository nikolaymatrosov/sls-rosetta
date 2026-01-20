# Random UUID for uploads bucket (unique across Yandex Cloud)
resource "random_uuid" "upload_bucket_name" {
}

# Bucket for image uploads and thumbnails
resource "yandex_storage_bucket" "for_uploads" {
  access_key = yandex_iam_service_account_static_access_key.sa_storage_editor.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa_storage_editor.secret_key
  bucket     = random_uuid.upload_bucket_name.result

  depends_on = [
    yandex_iam_service_account.sa_storage_editor,
    yandex_iam_service_account_static_access_key.sa_storage_editor,
    yandex_resourcemanager_folder_iam_member.sa_storage_editor
  ]
}
