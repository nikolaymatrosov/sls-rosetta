# Archive function code
data "archive_file" "function_code" {
  type        = "zip"
  source_dir  = "../function"
  output_path = "./storage-function.zip"
  excludes = [
    "__pycache__",
    "*.pyc",
    ".pytest_cache"
  ]
}

# Storage Handler Function
resource "yandex_function" "storage_handler" {
  name               = "storage-handler-python"
  description        = "Processes Object Storage events to create image thumbnails"
  user_hash          = data.archive_file.function_code.output_sha256
  runtime            = "python312"
  entrypoint         = "handler.handler"
  memory             = 256
  execution_timeout  = "30"
  service_account_id = yandex_iam_service_account.sa_storage_editor.id

  environment = {
    AWS_ACCESS_KEY_ID     = yandex_iam_service_account_static_access_key.sa_storage_editor.access_key
    AWS_SECRET_ACCESS_KEY = yandex_iam_service_account_static_access_key.sa_storage_editor.secret_key
  }

  content {
    zip_filename = data.archive_file.function_code.output_path
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.sa_storage_editor
  ]
}

# Storage Trigger
resource "yandex_function_trigger" "storage_trigger" {
  name        = "storage-trigger-python"
  description = "Trigger that invokes handler on new objects in uploads/ folder"
  folder_id   = var.folder_id

  object_storage {
    bucket_id    = yandex_storage_bucket.for_uploads.bucket
    prefix       = "uploads/"
    suffix       = ""
    create       = true
    update       = true
    batch_cutoff = 1
  }

  function {
    id                 = yandex_function.storage_handler.id
    service_account_id = yandex_iam_service_account.trigger_sa.id
    retry_attempts     = 3
    retry_interval     = 10
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.trigger_sa,
    yandex_storage_bucket.for_uploads
  ]
}
