resource "yandex_ydb_database_serverless" "db" {
  name      = "ydb-serverless-demo"
  folder_id = var.folder_id
  location_id = "ru-central1"
} 