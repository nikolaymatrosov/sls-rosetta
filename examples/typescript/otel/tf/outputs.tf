output "orders_function_id" {
  description = "ID of the orders service function"
  value       = yandex_function.otel_orders.id
}

output "payment_function_id" {
  description = "ID of the payment service function"
  value       = yandex_function.otel_payment.id
}

output "inventory_function_id" {
  description = "ID of the inventory service function"
  value       = yandex_function.otel_inventory.id
}

output "orders_url" {
  description = "URL to invoke the orders service function"
  value       = "https://functions.yandexcloud.net/${yandex_function.otel_orders.id}"
}

output "payment_url" {
  description = "URL to invoke the payment service function"
  value       = "https://functions.yandexcloud.net/${yandex_function.otel_payment.id}"
}

output "orders_queue_url" {
  description = "URL of the YMQ orders queue"
  value       = yandex_message_queue.orders_queue.id
}
