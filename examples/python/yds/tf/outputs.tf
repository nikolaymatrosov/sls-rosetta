output "producer_function_id" {
  description = "ID of the producer function"
  value       = yandex_function.producer.id
}

output "producer_function_url" {
  description = "URL of the producer function"
  value       = "https://functions.yandexcloud.net/${yandex_function.producer.id}"
}

output "consumer_function_id" {
  description = "ID of the consumer function"
  value       = yandex_function.consumer.id
}

output "ydb_database_id" {
  description = "ID of the YDB database"
  value       = yandex_ydb_database_serverless.yds_demo.id
}

output "ydb_database_endpoint" {
  description = "Endpoint of the YDB database"
  value       = yandex_ydb_database_serverless.yds_demo.ydb_full_endpoint
}

output "ydb_database_path" {
  description = "Path of the YDB database"
  value       = yandex_ydb_database_serverless.yds_demo.database_path
}

output "yds_topic_name" {
  description = "Name of the YDS topic"
  value       = yandex_ydb_topic.yds_demo_topic.name
}

output "yds_topic_path" {
  description = "Full path of the YDS topic"
  value       = "${yandex_ydb_database_serverless.yds_demo.database_path}/${yandex_ydb_topic.yds_demo_topic.name}"
}

output "trigger_id" {
  description = "ID of the YDS trigger"
  value       = yandex_function_trigger.yds_trigger.id
}

output "curl_example" {
  description = "Example curl command to test the producer"
  value       = <<-EOT
    curl -X POST https://functions.yandexcloud.net/${yandex_function.producer.id} \
      -H "Content-Type: application/json" \
      -d '{"message": "Hello YDS", "user_id": "user123", "action": "login"}'
  EOT
}
