resource "yandex_iam_service_account" "sa_serverless" {
  name      = "sa-serverless-test"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "sa_serverless" {
  for_each = toset([
    "serverless.functions.invoker",
  ])
  role      = each.value
  folder_id = var.folder_id
  members   = [
    "serviceAccount:${yandex_iam_service_account.sa_serverless.id}",
  ]
}