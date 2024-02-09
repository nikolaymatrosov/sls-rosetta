resource "yandex_ydb_database_serverless" "db" {
  name      = "test-ydb-serverless"
  folder_id = var.folder_id
}

resource "null_resource" "document-table" {
  provisioner "local-exec" {
    environment = {
      AWS_ACCESS_KEY_ID     = yandex_iam_service_account_static_access_key.db_admin.access_key
      AWS_SECRET_ACCESS_KEY = yandex_iam_service_account_static_access_key.db_admin.secret_key
    }
    command = <<EOF
          aws dynamodb create-table \
            --table-name demo \
            --attribute-definitions \
            AttributeName=id,AttributeType=N \
            AttributeName=key,AttributeType=S \
            AttributeName=value,AttributeType=S \
            --key-schema \
            AttributeName=id,KeyType=HASH \
            AttributeName=key,KeyType=RANGE \
            --endpoint ${yandex_ydb_database_serverless.db.document_api_endpoint}
EOF
  }
    depends_on = [
      yandex_ydb_database_serverless.db,
      yandex_iam_service_account_static_access_key.db_admin
    ]
}