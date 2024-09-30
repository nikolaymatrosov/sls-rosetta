resource "yandex_ydb_database_serverless" "example" {
  name      = "test-ydb-serverless"
  folder_id = var.folder_id
}

resource "yandex_ydb_table" "test_table" {
  path = "requests"
  connection_string = yandex_ydb_database_serverless.example.ydb_full_endpoint

  column {
    name = "id"
    type = "Utf8"
    not_null = true
  }
  column {
    name = "data"
    type = "JSONDocument"
    not_null = true
  }

  primary_key = ["id"]
}

resource "yandex_iam_service_account" "ydb_sa" {
  name      = "ydb-sa"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "ydb_service_account_binding" {
  for_each = toset([
    "ydb.admin",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.ydb_sa.id}",
  ]
  sleep_after = 5
}
