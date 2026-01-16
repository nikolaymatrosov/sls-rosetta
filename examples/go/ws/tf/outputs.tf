output "websocket_url" {
  description = "WebSocket connection URL"
  value       = "wss://${yandex_api_gateway.ws_gateway.domain}/ws"
}

output "gateway_id" {
  description = "API Gateway ID"
  value       = yandex_api_gateway.ws_gateway.id
}

output "function_id" {
  description = "Function ID"
  value       = yandex_function.ws_handler.id
}

output "database_id" {
  description = "YDB Database ID"
  value       = yandex_ydb_database_serverless.ws_database.id
}

output "database_endpoint" {
  description = "YDB Database endpoint"
  value       = yandex_ydb_database_serverless.ws_database.ydb_full_endpoint
}

output "service_account_id" {
  description = "Service Account ID"
  value       = yandex_iam_service_account.ws_sa.id
}

output "client_command" {
  description = "Command to run the WebSocket client"
  value       = "cd ../function/client && go run main.go -url wss://${yandex_api_gateway.ws_gateway.domain}/ws"
}
