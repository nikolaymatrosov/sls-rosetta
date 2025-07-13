resource "yandex_iam_service_account" "producer_sa" {
  name = "yds-producer-sa"
}

resource "yandex_iam_service_account_static_access_key" "producer_key" {
  service_account_id = yandex_iam_service_account.producer_sa.id
  description       = "Static access key for YDS producer function"
}

resource "yandex_resourcemanager_folder_iam_member" "producer_writer" {
  folder_id = var.folder_id
  role      = "ymq.writer"
  member    = "serviceAccount:${yandex_iam_service_account.producer_sa.id}"
}

variable "folder_id" {
  description = "Yandex Cloud folder ID"
  type        = string
} 