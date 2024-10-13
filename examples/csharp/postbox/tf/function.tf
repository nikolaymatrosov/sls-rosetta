data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../Postbox/Postbox"
  type        = "zip"
  excludes = [
    "bin/*",
    "obj/*",
  ]
}

resource "yandex_function" "postbox_function" {
  name              = "postbox-function"
  user_hash         = data.archive_file.function_files.output_sha256
  runtime           = "dotnet8"
  entrypoint        = "Postbox.Handler"
  memory            = "128"
  execution_timeout = "10"
  content {
    zip_filename = data.archive_file.function_files.output_path
  }
  environment = {
    AWS_SECRET_ACCESS_KEY = yandex_iam_service_account_static_access_key.postbox_sender_key.secret_key
    AWS_ACCESS_KEY_ID     = yandex_iam_service_account_static_access_key.postbox_sender_key.access_key
  }
}

// IAM binding for making function public
resource "yandex_function_iam_binding" "postbox_function_binding" {
  function_id = yandex_function.postbox_function.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}

