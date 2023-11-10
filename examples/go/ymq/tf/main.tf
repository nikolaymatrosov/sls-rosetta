resource "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../function"
  type        = "zip"
}

resource "yandex_function" "ymq_sender" {
  name               = "ymq-sender"
  user_hash          = archive_file.function_files.output_sha256
  runtime            = "golang121"
  entrypoint         = "sender.Sender"
  memory             = "128"
  execution_timeout  = "10"
  content {
    zip_filename = archive_file.function_files.output_path
  }
  service_account_id = yandex_iam_service_account.ymq_writer.id
  environment = {
    "YMQ_NAME" = yandex_message_queue.input_queue.name
    "AWS_ACCESS_KEY_ID" = yandex_iam_service_account_static_access_key.ymq_writer.access_key
    "AWS_SECRET_ACCESS_KEY" = yandex_iam_service_account_static_access_key.ymq_writer.secret_key
  }
}


// IAM binding for making function public
resource "yandex_function_iam_binding" "ymq_sender_binding" {
  function_id = yandex_function.ymq_sender.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}


resource "yandex_function" "ymq_receiver" {
  name               = "ymq-receiver"
  user_hash          = archive_file.function_files.output_sha256
  runtime            = "golang121"
  entrypoint         = "receiver.Receiver"
  memory             = "128"
  execution_timeout  = "10"
  content {
    zip_filename = archive_file.function_files.output_path
  }
  // This function also need ability to write to message queue
  // because that is the way it will return the result of execution
  service_account_id = yandex_iam_service_account.ymq_writer.id
  environment = {
    "YMQ_NAME" = yandex_message_queue.response_queue.name
    "AWS_ACCESS_KEY_ID" = yandex_iam_service_account_static_access_key.ymq_writer.access_key
    "AWS_SECRET_ACCESS_KEY" = yandex_iam_service_account_static_access_key.ymq_writer.secret_key
  }
}

resource "yandex_function_trigger" "ymq_trigger" {
  name        = "ymq-trigger"

  message_queue {
    queue_id = yandex_message_queue.input_queue.arn
    batch_cutoff = "5"
    batch_size = "5"
    service_account_id = yandex_iam_service_account.trigger_sa.id
  }
  function {
    id = yandex_function.ymq_receiver.id
    service_account_id = yandex_iam_service_account.trigger_sa.id
  }
}