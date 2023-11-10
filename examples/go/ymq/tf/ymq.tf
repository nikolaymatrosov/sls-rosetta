resource "yandex_message_queue" "input_queue" {
  name                       = "input_queue"
  visibility_timeout_seconds = 600
  receive_wait_time_seconds  = 20
  message_retention_seconds  = 1209600
  redrive_policy             = jsonencode({
    deadLetterTargetArn = yandex_message_queue.example_deadletter_queue.arn
    maxReceiveCount     = 3
  })
  access_key = yandex_iam_service_account_static_access_key.sa_ymq_creator.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa_ymq_creator.secret_key
}

resource "yandex_message_queue" "example_deadletter_queue" {
  name       = "ymq_deadletter_example"
  access_key = yandex_iam_service_account_static_access_key.sa_ymq_creator.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa_ymq_creator.secret_key
}

resource "yandex_message_queue" "response_queue" {
  name                       = "response_queue"
  visibility_timeout_seconds = 600
  message_retention_seconds  = 1209600
  access_key                 = yandex_iam_service_account_static_access_key.sa_ymq_creator.access_key
  secret_key                 = yandex_iam_service_account_static_access_key.sa_ymq_creator.secret_key
}