resource "yandex_iam_service_account" "storage-admin" {
  folder_id   = var.folder_id
  description = "Service account for storage admin"
  name        = "storage-admin"
}

resource "yandex_resourcemanager_folder_iam_binding" "storage-admin" {
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.storage-admin.id}"
  ]
  role = "storage.admin"
}

resource "yandex_iam_service_account_static_access_key" "sa-static-key" {
  service_account_id = yandex_iam_service_account.storage-admin.id
  description        = "Static access key for storage admin"
}

resource "yandex_iam_service_account" "url-shortener" {
  folder_id   = var.folder_id
  description = "Service account for URL shortener"
  name        = "url-shortener"
}

locals {
  sa_roles = [
    "storage.viewer",
    "ydb.admin",
    "functions.functionInvoker",
  ]
}
resource "yandex_resourcemanager_folder_iam_binding" "url-shortener" {
  for_each  = toset(local.sa_roles)
  folder_id = var.folder_id
  role      = each.key
  members   = [
    "serviceAccount:${yandex_iam_service_account.url-shortener.id}"
  ]
}

resource "yandex_storage_bucket" "url-shortener-bucket" {
  access_key = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
  bucket     = "url-shortener-bucket"
}

resource "yandex_storage_object" "index_html" {
  bucket = yandex_storage_bucket.url-shortener-bucket.bucket
  key    = "index.html"
  source = "./index.html"
  content_type = "text/html"
}

resource "yandex_ydb_database_serverless" "url-shortener-db" {
  name      = "url-shortener-db"
  folder_id = var.folder_id
}

resource "yandex_ydb_table" "test_table" {
  path              = "links"
  connection_string = yandex_ydb_database_serverless.url-shortener-db.ydb_full_endpoint

  column {
    name = "id"
    type = "Utf8"
  }
  column {
    name = "link"
    type = "Utf8"
  }
  primary_key = ["id"]

}

data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "function"
  type        = "zip"
}

resource "yandex_function" "url-shortener" {
  name              = "url-shortener"
  user_hash         = data.archive_file.function_files.output_sha256
  runtime           = "python312"
  entrypoint        = "index.handler"
  memory            = "256"
  execution_timeout = "10"
  content {
    zip_filename = data.archive_file.function_files.output_path
  }
  environment = {
    USE_METADATA_CREDENTIALS = 1
    endpoint                 = "grpcs://ydb.serverless.yandexcloud.net:2135"
    database                 = yandex_ydb_database_serverless.url-shortener-db.database_path
  }
  service_account_id = yandex_iam_service_account.url-shortener.id
}

resource "yandex_api_gateway" "url-shortener-api" {
  name        = "url-shortener-api"
  description = "API for URL shortener"

  spec = templatefile("./api-gateway.yaml", {
    function_id        = yandex_function.url-shortener.id,
    bucket             = yandex_storage_bucket.url-shortener-bucket.bucket,
    service_account_id = yandex_iam_service_account.url-shortener.id,
  })
}

output "url_shortener_api_url" {
  value = "https://${yandex_api_gateway.url-shortener-api.domain}"
}