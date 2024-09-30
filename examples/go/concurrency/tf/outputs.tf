output "simple" {
  value = yandex_function.simple.id
}

output "long" {
  value = yandex_function.long.id
}

output "ydb" {
  value = yandex_function.ydb.id
}

output "ydb_database_url" {
  value = yandex_ydb_database_serverless.example.ydb_full_endpoint
}
