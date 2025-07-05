output "function_id" {
  value = yandex_function.ydb_function.id
}

output "function_url" {
  value = "https://functions.yandexcloud.net/${yandex_function.ydb_function.id}"
}

output "ydb_database_path" {
  value = yandex_ydb_database_serverless.db.database_path
}

output "ydb_endpoint" {
  value = yandex_ydb_database_serverless.db.ydb_api_endpoint
} 