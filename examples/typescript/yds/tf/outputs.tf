output "producer_function_id" {
  value = yandex_function.producer.id
}

output "yds_topic_name" {
  value = yandex_message_queue_topic.yds_topic.name
} 