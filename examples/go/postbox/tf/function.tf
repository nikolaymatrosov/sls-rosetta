data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../function"
  type        = "zip"
}

resource "yandex_function" "postbox_aws" {
  name              = "postbox-aws"
  user_hash         = data.archive_file.function_files.output_sha256
  runtime           = "golang121"
  entrypoint        = "index.AwsHandler"
  memory            = "128"
  execution_timeout = "10"
  content {
    zip_filename = data.archive_file.function_files.output_path
  }
  service_account_id = yandex_iam_service_account.postbox_sender.id
  environment = {
    AWS_SECRET_ACCESS_KEY = yandex_iam_service_account_static_access_key.postbox_sender_key.secret_key
    AWS_ACCESS_KEY_ID     = yandex_iam_service_account_static_access_key.postbox_sender_key.access_key
  }
}

resource "yandex_function" "postbox_yc" {
  name              = "postbox-yc"
  user_hash         = data.archive_file.function_files.output_sha256
  runtime           = "golang121"
  entrypoint        = "index.YcHandler"
  memory            = "128"
  execution_timeout = "10"
  content {
    zip_filename = data.archive_file.function_files.output_path
  }
  service_account_id = yandex_iam_service_account.postbox_sender.id
}

locals {
  functions = {
    aws = yandex_function.postbox_aws,
    yc  = yandex_function.postbox_yc,
  }
}

// IAM binding for making function public
resource "yandex_function_iam_binding" "postbox_function_binding" {
  for_each    = local.functions
  function_id = each.value.id
  role        = "functions.functionInvoker"
  members = ["system:allUsers"]
}
