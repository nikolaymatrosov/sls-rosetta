resource "yandex_ydb_database_serverless" "default" {
  name = "default"
  location_id = "ru-central1"
}


resource "yandex_ydb_topic" "input-topic" {
  database_endpoint = "${yandex_ydb_database_serverless.default.ydb_full_endpoint}"
  name = "yds-topic"

  supported_codecs = ["raw", "gzip"]
  partitions_count = 1
  retention_period_ms = 2000000
}