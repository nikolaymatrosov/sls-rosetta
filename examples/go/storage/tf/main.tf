locals {
  archive_output_path = "../environment/function.zip"
}

data "archive_file" "function_code" {
  output_path = local.archive_output_path
  source_dir  = "../build"
  type        = "zip"
}

resource "yandex_storage_object" "function_code" {
  access_key = yandex_iam_service_account_static_access_key.sa_storage_editor.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa_storage_editor.secret_key

  bucket      = yandex_storage_bucket.for-deploy.bucket
  key         = "function.zip"
  source      = local.archive_output_path
  source_hash = data.archive_file.function_code.output_sha
  depends_on  = [
    data.archive_file.function_code,
  ]
}

resource "yandex_function" "storage-handler" {
  name              = "storage-handler"
  user_hash         = data.archive_file.function_code.output_sha
  runtime           = "golang123"
  entrypoint        = "handler.Handler"
  memory            = "128"
  execution_timeout = "10"
  package {
    bucket_name = yandex_storage_bucket.for-deploy.bucket
    object_name = "function.zip"
    #    sha_256     = archive_file.function_code.output_sha256
  }
  service_account_id = yandex_iam_service_account.sa_storage_editor.id
  environment        = {
    # The trigger will provide the name of the bucket and object key, but not actual content of the object
    # So we need to get the content of the object ourselves
    "AWS_ACCESS_KEY_ID"     = yandex_iam_service_account_static_access_key.sa_storage_editor.access_key
    "AWS_SECRET_ACCESS_KEY" = yandex_iam_service_account_static_access_key.sa_storage_editor.secret_key
  }
  depends_on = [
    yandex_storage_object.function_code
  ]
}

resource "yandex_function_trigger" "storage-trigger" {
  name = "storage-trigger"

  object_storage {
    bucket_id    = yandex_storage_bucket.for-uploads.bucket
    prefix       = "uploads/"
    batch_cutoff = 1
    create       = true
    update       = true
  }
  function {
    id                 = yandex_function.storage-handler.id
    service_account_id = yandex_iam_service_account.trigger_sa.id
  }
}