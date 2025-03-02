resource "yandex_api_gateway" "api_gateway" {
  name = "api-gateway"
  spec = templatefile("./api-gateway.yaml", {
    test_function_id   = yandex_function.test_function.id
    route_function_id  = yandex_function.route_function.id
    sa_id = yandex_iam_service_account.sa_serverless.id
  })
}