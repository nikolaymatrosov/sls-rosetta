# Create service account for the function
resource "yandex_iam_service_account" "ws_sa" {
  name        = var.service_account_name
  description = "Service account for WebSocket Go example"
  folder_id   = var.folder_id
}

# Grant YDB editor role to the service account
resource "yandex_resourcemanager_folder_iam_member" "sa_ydb_editor" {
  folder_id = var.folder_id
  role      = "ydb.editor"
  member    = "serviceAccount:${yandex_iam_service_account.ws_sa.id}"
}

# Grant WebSocket writer role to the service account
resource "yandex_resourcemanager_folder_iam_member" "sa_websocket_writer" {
  folder_id = var.folder_id
  role      = "api-gateway.websocketWriter"
  member    = "serviceAccount:${yandex_iam_service_account.ws_sa.id}"
}

# Grant WebSocket broadcaster role to the service account
resource "yandex_resourcemanager_folder_iam_member" "sa_websocket_broadcaster" {
  folder_id = var.folder_id
  role      = "api-gateway.websocketBroadcaster"
  member    = "serviceAccount:${yandex_iam_service_account.ws_sa.id}"
}

# Grant function invoker role to the service account (for triggers)
resource "yandex_resourcemanager_folder_iam_member" "sa_functions_invoker" {
  folder_id = var.folder_id
  role      = "serverless.functions.invoker"
  member    = "serviceAccount:${yandex_iam_service_account.ws_sa.id}"
}

# Grant YDS admin role to the service account (for managing topics and triggers)
resource "yandex_resourcemanager_folder_iam_member" "sa_yds_admin" {
  folder_id = var.folder_id
  role      = "yds.admin"
  member    = "serviceAccount:${yandex_iam_service_account.ws_sa.id}"
}
