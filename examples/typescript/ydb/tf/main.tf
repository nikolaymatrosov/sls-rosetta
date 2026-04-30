# Build TypeScript code
resource "null_resource" "build_typescript" {
  provisioner "local-exec" {
    command     = "cd ../function && npm install && npm run build"
    working_dir = path.module
  }

  triggers = {
    main_hash    = filesha256("../function/main.ts")
    package_hash = filesha256("../function/package.json")
    tsconfig_hash = filesha256("../function/tsconfig.json")
  }
}

# Archive built function code
data "archive_file" "function_files" {
  type        = "zip"
  source_dir  = "../dist"
  output_path = "./function.zip"
  excludes = [
    "node_modules",
    "*.ts",
    "tsconfig.json"
  ]

  depends_on = [null_resource.build_typescript]
}

resource "yandex_function" "ydb_function" {
  name               = "ydb-demo-ts"
  description        = "YDB query with transaction statistics (TypeScript)"
  user_hash          = data.archive_file.function_files.output_sha256
  runtime            = "nodejs22"
  entrypoint         = "handler.handler"
  memory             = 128
  execution_timeout  = "10"
  service_account_id = yandex_iam_service_account.function_sa.id

  content {
    zip_filename = data.archive_file.function_files.output_path
  }

  environment = {
    YDB_DATABASE = yandex_ydb_database_serverless.db.database_path
    YDB_ENDPOINT = yandex_ydb_database_serverless.db.ydb_api_endpoint
  }

  depends_on = [
    yandex_ydb_database_serverless.db,
    null_resource.run_migrations,
    yandex_iam_service_account.function_sa,
    yandex_resourcemanager_folder_iam_binding.function_sa,
  ]
}

# IAM binding for making function public
resource "yandex_function_iam_binding" "function_binding" {
  function_id = yandex_function.ydb_function.id
  role        = "functions.functionInvoker"
  members     = ["system:allUsers"]
}
