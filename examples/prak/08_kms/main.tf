resource "yandex_iam_service_account" "kms" {
  folder_id   = var.folder_id
  description = "Service account for KMS"
  name        = "kms"
}

resource "yandex_resourcemanager_folder_iam_binding" "kms-roles" {
  folder_id = var.folder_id
  role      = "kms.admin"
  members   = [
    "serviceAccount:${yandex_iam_service_account.kms.id}"
  ]
}

resource "yandex_kms_symmetric_key" "praktikum" {
  folder_id = var.folder_id
  name      = "praktikum"
  description = "Key for praktikum"
  default_algorithm = "AES_128"
  rotation_period   = "8760h" // equal to 1 year
}

