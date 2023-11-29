# Bucket names are unique across all Yandex.Cloud users.
# We can not just use uuid function because it will be different for each run
# So we need to use random_uuid resource
resource "random_uuid" "upload-bucket-name" {
}

resource "yandex_storage_bucket" "for-uploads" {
  access_key = yandex_iam_service_account_static_access_key.sa_storage_editor.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa_storage_editor.secret_key
  bucket     = random_uuid.upload-bucket-name.result
  depends_on = [
    yandex_iam_service_account.sa_storage_editor
  ]
}

resource "random_uuid" "deploy-bucket-name" {}

resource "yandex_storage_bucket" "for-deploy" {
  access_key = yandex_iam_service_account_static_access_key.sa_storage_editor.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa_storage_editor.secret_key
  bucket     = random_uuid.deploy-bucket-name.result
  depends_on = [
    yandex_iam_service_account.sa_storage_editor
  ]
}