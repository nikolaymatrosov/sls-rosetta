output "ymq_id" {
  value = yandex_message_queue.response_queue.id
}

output "ymq_arn" {
  value = yandex_message_queue.response_queue.arn
}

output "ymq_name" {
  value = yandex_message_queue.response_queue.name
}

output "ymq_reader_access_key" {
  value     = yandex_iam_service_account_static_access_key.ymq_reader.access_key
  sensitive = true
}

output "ymq_reader_secret_key" {
  value     = yandex_iam_service_account_static_access_key.ymq_reader.secret_key
  sensitive = true
}

output "sender_function_id" {
  value = yandex_function.ymq_sender.id
}