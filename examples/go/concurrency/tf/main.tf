locals {
  concurrency = 16
}

resource "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../function"
  type        = "zip"
}

resource "yandex_function" "simple" {
  name              = "con-simple"
  user_hash         = archive_file.function_files.output_sha256
  runtime           = "golang121"
  entrypoint        = "simple.Simple"
  memory            = "128"
  execution_timeout = "60"
  concurrency       = local.concurrency
  content {
    zip_filename = archive_file.function_files.output_path
  }
  environment = {
    "CONCURRENCY" = local.concurrency
  }
}

resource "yandex_function" "long" {
  name              = "con-long"
  user_hash         = archive_file.function_files.output_sha256
  runtime           = "golang121"
  entrypoint        = "long.Long"
  memory            = "128"
  execution_timeout = "10"
  concurrency       = local.concurrency
  content {
    zip_filename = archive_file.function_files.output_path
  }
  environment = {
    "CONCURRENCY" = local.concurrency
  }
}

resource "yandex_function" "ydb" {
  name               = "con-ydb"
  user_hash          = archive_file.function_files.output_sha256
  runtime            = "golang121"
  entrypoint         = "ydb.YdbHandler"
  memory             = "128"
  execution_timeout  = "10"
  concurrency        = local.concurrency
  service_account_id = yandex_iam_service_account.ydb_sa.id

  content {
    zip_filename = archive_file.function_files.output_path
  }
  environment = {
    "CONCURRENCY"              = local.concurrency
    "YDB_DSN"                  = yandex_ydb_database_serverless.example.ydb_full_endpoint
    "YDB_METADATA_CREDENTIALS" = "1"
  }
}


locals {
  function_ids = {
    simple : yandex_function.simple.id,
    long : yandex_function.long.id,
    ydb : yandex_function.ydb.id
  }
}

// IAM binding for making function public
resource "yandex_function_iam_binding" "test_function_binding" {
  for_each    = local.function_ids
  function_id = each.value
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}



