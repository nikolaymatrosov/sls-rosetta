resource "yandex_ydb_database_serverless" "db" {
  name      = "test-ydb-serverless"
  folder_id = var.folder_id
}

resource "aws_dynamodb_table" "test" {
  depends_on = [
    yandex_resourcemanager_folder_iam_binding.db_admin,
    yandex_iam_service_account_static_access_key.db_admin,
    yandex_ydb_database_serverless.db
  ]
  name         = "demo"
  billing_mode = "PAY_PER_REQUEST" # только такой billing_mode поддерживается у нас и его нужно явно указывать.

  hash_key  = "id"
  range_key = "key"

  attribute {
    name = "id"
    type = "N"
  }

  attribute {
    name = "key"
    type = "S"
  }
}
