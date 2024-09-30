resource "yandex_kms_symmetric_key" "kms_key" {
  name              = "k8s-kms-key"
  description       = "K8S KMS symetric key"
  default_algorithm = "AES_256"
  rotation_period   = "8760h"
}
