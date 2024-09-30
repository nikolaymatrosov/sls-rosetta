resource "yandex_logging_group" "k8s_log_group" {
  name      = "k8s-logging-group"
  folder_id = var.folder_id
}
