# YDB Serverless Database
resource "yandex_ydb_database_serverless" "yds_demo" {
  name        = "yds-demo-db-typescript"
  location_id = "ru-central1"

  serverless_database {
    storage_size_limit = 5
  }

  sleep_after = 5
}

# YDB Topic (Data Stream)
resource "yandex_ydb_topic" "yds_demo_topic" {
  database_endpoint = yandex_ydb_database_serverless.yds_demo.ydb_full_endpoint
  name              = "yds-demo-topic"

  supported_codecs = ["raw"]

  partitions_count       = 1
  retention_period_hours = 24
}
