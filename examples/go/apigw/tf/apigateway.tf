resource "yandex_api_gateway" "api_gateway" {
  name = "api-gateway"
  spec = <<-EOT
    openapi: "3.0.0"
    info:
      version: 1.0.0
      title: Test API
    paths:
      /demo:
        post:
          operationId: demo
          x-yc-apigateway-integration:
            type: cloud_functions
            function_id: ${yandex_function.test_function.id}
            service_account_id: ${yandex_iam_service_account.sa_serverless.id}
            payload_format_version: "1.0"
  EOT
}