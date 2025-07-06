output "producer_function_id" {
  value = yandex_function.producer_function.id
}

output "producer_function_url" {
  value = "https://functions.yandexcloud.net/${yandex_function.producer_function.id}"
}

output "consumer_function_id" {
  value = yandex_function.consumer_function.id
}

output "yds_topic_id" {
  value = yandex_ydb_topic.main_topic.id
}

output "yds_topic_name" {
  value = yandex_ydb_topic.main_topic.name
}

output "yds_database_id" {
  value = yandex_ydb_database_serverless.yds_db.id
}

output "yds_database_path" {
  value = yandex_ydb_database_serverless.yds_db.database_path
}

output "trigger_id" {
  value = yandex_function_trigger.yds_trigger.id
} 