# Archive function code
data "archive_file" "producer_function" {
  type        = "zip"
  source_dir  = "../function"
  output_path = "./producer-function.zip"
  excludes = [
    "__pycache__",
    "*.pyc",
    ".pytest_cache"
  ]
}

data "archive_file" "consumer_function" {
  type        = "zip"
  source_dir  = "../function"
  output_path = "./consumer-function.zip"
  excludes = [
    "__pycache__",
    "*.pyc",
    ".pytest_cache"
  ]
}

# Producer Function
resource "yandex_function" "producer" {
  name               = "yds-producer-python"
  description        = "Producer function that writes messages to YDS topic"
  user_hash          = data.archive_file.producer_function.output_sha256
  runtime            = "python312"
  entrypoint         = "producer.producer_handler"
  memory             = 128
  execution_timeout  = "10"
  service_account_id = yandex_iam_service_account.producer_sa.id

  environment = {
    YDB_ENDPOINT   = yandex_ydb_database_serverless.yds_demo.ydb_full_endpoint
    YDS_TOPIC_PATH = "${yandex_ydb_database_serverless.yds_demo.database_path}/${yandex_ydb_topic.yds_demo_topic.name}"
  }

  content {
    zip_filename = data.archive_file.producer_function.output_path
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.producer_roles
  ]
}

# Consumer Function
resource "yandex_function" "consumer" {
  name               = "yds-consumer-python"
  description        = "Consumer function that processes messages from YDS topic"
  user_hash          = data.archive_file.consumer_function.output_sha256
  runtime            = "python312"
  entrypoint         = "consumer.consumer_handler"
  memory             = 128
  execution_timeout  = "10"
  service_account_id = yandex_iam_service_account.consumer_sa.id

  content {
    zip_filename = data.archive_file.consumer_function.output_path
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.consumer_roles
  ]
}

# YDS Trigger
resource "yandex_function_trigger" "yds_trigger" {
  name        = "yds-trigger-python"
  description = "Trigger that invokes consumer function on YDS messages"
  folder_id   = var.folder_id

  function {
    id                 = yandex_function.consumer.id
    service_account_id = yandex_iam_service_account.trigger_sa.id
    retry_attempts     = 3
    retry_interval     = 10
  }

  data_streams {
    stream_name = yandex_ydb_topic.yds_demo_topic.name
    database    = yandex_ydb_database_serverless.yds_demo.database_path
    service_account_id = yandex_iam_service_account.trigger_sa.id
    batch_cutoff = 5
    batch_size   = 10
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_member.trigger_roles,
    yandex_ydb_topic.yds_demo_topic
  ]
}
