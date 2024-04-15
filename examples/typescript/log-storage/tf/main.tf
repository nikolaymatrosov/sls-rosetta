resource "null_resource" "build_typescript" {
  provisioner "local-exec" {
    command = "cd ../function && npm run build"
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
resource "yandex_function" "log_generator_function" {
  name               = "logs-generator"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs18"
  entrypoint         = "main.handler"
  memory             = "128"
  execution_timeout  = "10"
  content {
    zip_filename = data.archive_file.function_files.output_path
  }
}

// IAM binding for making function public
resource "yandex_function_iam_binding" "test_function_binding" {
  function_id = yandex_function.log_generator_function.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}



