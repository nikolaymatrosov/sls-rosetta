resource "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../function"
  type        = "zip"
}

resource "yandex_function" "sdk" {
  name               = "sdk-demo"
  user_hash          = archive_file.function_files.output_sha256
  runtime            = "golang121"
  entrypoint         = "index.Handler"
  memory             = "128"
  execution_timeout  = "10"
  content {
    zip_filename = archive_file.function_files.output_path
  }
}

// IAM binding for making function public
resource "yandex_function_iam_binding" "test_function_binding" {
  function_id = yandex_function.sdk.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}



