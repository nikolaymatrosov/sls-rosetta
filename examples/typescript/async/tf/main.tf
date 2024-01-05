resource "null_resource" "build_typescript" {
  provisioner "local-exec" {
    command = "cd ../function && npm i && npm run build"
  }
  triggers = {
    always_run = timestamp()
  }
}

data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../dist"
  type        = "zip"
  depends_on  = [
    null_resource.build_typescript
  ]
}

resource "yandex_function" "async_function" {
  name               = "async-function"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs18"
  entrypoint         = "main.handler"
  memory             = "128"
  execution_timeout  = "10"

  content {
    zip_filename = data.archive_file.function_files.output_path
  }

  async_invocation {
    retries_count = "3"
    service_account_id = yandex_iam_service_account.function_invoker.id
    ymq_failure_target {
      service_account_id = yandex_iam_service_account.ymq_writer.id
      arn = yandex_message_queue.failed_queue.arn
    }
    ymq_success_target {
      service_account_id = yandex_iam_service_account.ymq_writer.id
      arn = yandex_message_queue.success_queue.arn
    }
  }

  depends_on = [
    yandex_resourcemanager_folder_iam_binding.function_invoker
  ]
}

// IAM binding for making function public
resource "yandex_function_iam_binding" "test_function_binding" {
  function_id = yandex_function.async_function.id
#  role        = "serverless.functions.invoker"
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}



