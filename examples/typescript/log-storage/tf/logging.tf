resource "yandex_ydb_database_serverless" "log_db" {
  name      = "logging_db"
  folder_id = var.folder_id
}

resource "yandex_ydb_topic" "log_topic" {
  database_endpoint = yandex_ydb_database_serverless.log_db.endpoint
  name              = "function-log-topic"

  supported_codecs           = ["raw", "gzip"]
  partitions_count           = 1
  retention_period_ms        = 2000000
  partition_write_speed_kbps = 128
}

resource "yandex_logging_group" "function_logs" {
  name        = "function-logs"
  folder_id   = var.folder_id
  data_stream = "${yandex_ydb_database_serverless.log_db.database_path}/${yandex_ydb_topic.log_topic.name}"
}

resource "yandex_datatransfer_endpoint" "topic_source" {
  name      = "logging-topic-source"
  folder_id = var.folder_id
  settings {
    yds_source {
      database           = yandex_ydb_database_serverless.log_db.database_path
      endpoint           = yandex_ydb_database_serverless.log_db.endpoint
      stream             = yandex_ydb_topic.log_topic.name
      service_account_id = yandex_iam_service_account.logging_transfer.id
    }
  }
}

resource "yandex_datatransfer_transfer" "storage_dest" {
  name      = "storage-dest"
  folder_id = var.folder_id
  source_id = yandex_datatransfer_endpoint.topic_source.id
  destination {
    storage_ {
      database           = yandex_ydb_database_serverless.log_db.database_path
      endpoint           = yandex_ydb_database_serverless.log_db.endpoint
      stream             = yandex_ydb_topic.log_topic.name
      service_account_id = yandex_iam_service_account.logging_transfer.id
    }
  }
}
```