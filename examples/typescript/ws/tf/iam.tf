resource "yandex_iam_service_account" "ws_function_sa" {
  name      = "ws-broadcast-function-sa"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "ws_function_sa" {
  for_each = toset([
    "ydb.editor",
    "api-gateway.websocketWriter",
    "api-gateway.websocketBroadcaster",
    "serverless.functions.invoker",
    "yds.admin",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.ws_function_sa.id}",
  ]
}
