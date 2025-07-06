# YDB Database for Data Streams
resource "yandex_ydb_database_serverless" "yds_db" {
  name        = "yds-demo-db"
  folder_id   = var.folder_id
  location_id = "ru-central1"

  sleep_after = 5
}

# YDB Topic for data ingestion (Data Streams)
resource "yandex_ydb_topic" "main_topic" {
  name = "yds-demo-topic"
  supported_codecs = [
    "raw",
  ]
  database_endpoint      = yandex_ydb_database_serverless.yds_db.ydb_full_endpoint
  description            = "Demo topic for YDS trigger example"
  partitions_count       = 1
  retention_period_hours = 24
} 