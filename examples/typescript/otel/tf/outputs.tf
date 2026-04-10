output "caller_function_id" {
  description = "ID of the caller function"
  value       = yandex_function.otel_caller.id
}

output "callee_function_id" {
  description = "ID of the callee function"
  value       = yandex_function.otel_callee.id
}

output "caller_url" {
  description = "URL to invoke the caller function"
  value       = "https://functions.yandexcloud.net/${yandex_function.otel_caller.id}"
}

output "callee_url" {
  description = "URL to invoke the callee function"
  value       = "https://functions.yandexcloud.net/${yandex_function.otel_callee.id}"
}
