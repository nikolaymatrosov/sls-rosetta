resource "yandex_api_gateway" "ws_gateway" {
  name = "websocket-broadcast-gateway"
  spec = templatefile("./api-gateway.yaml", {
    function_id    = yandex_function.ws_handler.id
    function_v2_id = yandex_function.ws_handler_v2.id
    sa_id          = yandex_iam_service_account.ws_function_sa.id
  })

  depends_on = [
    yandex_function.ws_handler,
    yandex_function.ws_handler_v2,
  ]
}
