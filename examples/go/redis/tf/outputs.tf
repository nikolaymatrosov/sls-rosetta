output "function_id" {
  value = {for k, v in yandex_function.redis_function : k => v.id}
}
