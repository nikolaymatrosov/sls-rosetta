# Create API Gateway for WebSocket
resource "yandex_api_gateway" "ws_gateway" {
  name = var.gateway_name

  spec = templatefile("${path.module}/api-gateway.yaml", {
    function_id = yandex_function.ws_handler.id
    sa_id       = yandex_iam_service_account.ws_sa.id
  })

  depends_on = [
    yandex_function.ws_handler,
    yandex_iam_service_account.ws_sa,
  ]
}
