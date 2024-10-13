data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../function"
  type        = "zip"
}

resource "yandex_function" "postbox_function" {
  name              = "postbox-function"
  user_hash         = data.archive_file.function_files.output_sha256
  runtime           = "golang121"
  entrypoint        = "index.Handler"
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
