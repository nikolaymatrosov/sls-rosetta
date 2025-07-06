# Archive the function code
resource "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../function"
  type        = "zip"
}

locals {
  runtime = "golang123"
}

# Producer function - writes data to YDS stream
resource "yandex_function" "producer_function" {
  name               = "yds-producer"
  user_hash          = archive_file.function_files.output_sha256
  runtime            = local.runtime
  entrypoint         = "main.ProducerHandler"
  memory             = "128"
  execution_timeout  = "10"
  service_account_id = yandex_iam_service_account.producer_sa.id

  content {
    zip_filename = archive_file.function_files.output_path
  }

  environment = {
    YDB_ENDPOINT  = yandex_ydb_database_serverless.yds_db.ydb_full_endpoint
    YDS_TOPIC_ID  = yandex_ydb_topic.main_topic.name
  }

  depends_on = [
    yandex_ydb_database_serverless.yds_db,
    yandex_ydb_topic.main_topic,
    yandex_iam_service_account.producer_sa,
    yandex_resourcemanager_folder_iam_binding.producer_sa,
  ]
}

# Consumer function - triggered by YDS events
resource "yandex_function" "consumer_function" {
  name               = "yds-consumer"
  user_hash          = archive_file.function_files.output_sha256
  runtime            = local.runtime
  entrypoint         = "main.ConsumerHandler"
  memory             = "128"
  execution_timeout  = "10"
  service_account_id = yandex_iam_service_account.consumer_sa.id

  content {
    zip_filename = archive_file.function_files.output_path
  }

  depends_on = [
    yandex_ydb_database_serverless.yds_db,
    yandex_ydb_topic.main_topic,
    yandex_iam_service_account.consumer_sa,
    yandex_resourcemanager_folder_iam_binding.consumer_sa,
  ]
}

# YDS Trigger - links topic to consumer function
resource "yandex_function_trigger" "yds_trigger" {
  name = "yds-trigger"

  data_streams {
    database           = yandex_ydb_database_serverless.yds_db.database_path
    stream_name        = yandex_ydb_topic.main_topic.name
    batch_cutoff       = "5"
    batch_size         = "10"
    service_account_id = yandex_iam_service_account.trigger_sa.id
  }

  function {
    id                 = yandex_function.consumer_function.id
    service_account_id = yandex_iam_service_account.trigger_sa.id
  }

  depends_on = [
    yandex_ydb_database_serverless.yds_db,
    yandex_ydb_topic.main_topic,
    yandex_function.consumer_function,
    yandex_iam_service_account.trigger_sa,
    yandex_resourcemanager_folder_iam_binding.trigger_sa,
  ]
}

# IAM binding for making producer function public
resource "yandex_function_iam_binding" "producer_binding" {
  function_id = yandex_function.producer_function.id
  role        = "functions.functionInvoker"
  members = ["system:allUsers"]
} 