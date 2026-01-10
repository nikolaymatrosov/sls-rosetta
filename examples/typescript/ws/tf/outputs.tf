output "api_gateway_id" {
  value = yandex_api_gateway.ws_gateway.id
}

output "websocket_url" {
  value = "wss://${yandex_api_gateway.ws_gateway.domain}/ws"
}

output "ydb_database_path" {
  value = yandex_ydb_database_serverless.ws_db.database_path
}

output "ydb_endpoint" {
  value = yandex_ydb_database_serverless.ws_db.ydb_api_endpoint
}

output "function_id" {
  value = yandex_function.ws_handler.id
}

output "migrate" {
  value = "grpcs://${yandex_ydb_database_serverless.ws_db.ydb_api_endpoint}${yandex_ydb_database_serverless.ws_db.database_path}"
}
