resource "yandex_message_queue" "orders_queue" {
  name                       = "otel-orders-queue"
  visibility_timeout_seconds = 600
  receive_wait_time_seconds  = 20
  message_retention_seconds  = 1209600

  redrive_policy = jsonencode({
    deadLetterTargetArn = yandex_message_queue.orders_dlq.arn
    maxReceiveCount     = 3
  })

  access_key = yandex_iam_service_account_static_access_key.otel_mq_admin_sa.access_key
  secret_key = yandex_iam_service_account_static_access_key.otel_mq_admin_sa.secret_key

  depends_on = [yandex_message_queue.orders_dlq]
}

resource "yandex_message_queue" "orders_dlq" {
  name = "otel-orders-dlq"

  access_key = yandex_iam_service_account_static_access_key.otel_mq_admin_sa.access_key
  secret_key = yandex_iam_service_account_static_access_key.otel_mq_admin_sa.secret_key

  depends_on = [
    yandex_iam_service_account_static_access_key.otel_mq_admin_sa,
  ]
}
