data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../function"
  type        = "zip"
  excludes = [
    "pom.xml",
  ]
}

resource "yandex_function" "hello_function" {
  name              = "hello-world"
  user_hash         = data.archive_file.function_files.output_sha256
  runtime           = "java21"
  entrypoint        = "ru.nikolaymatrosov.Handler"
  memory            = "128"
  execution_timeout = "10"
  content {
    zip_filename = data.archive_file.function_files.output_path
  }
}

// IAM binding for making function public
resource "yandex_function_iam_binding" "test_function_binding" {
  function_id = yandex_function.hello_function.id
  role        = "functions.functionInvoker"
  members = ["system:allUsers"]
}



