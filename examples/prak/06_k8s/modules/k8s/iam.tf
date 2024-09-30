locals {
  service_account = {
    master           = "k8s-master-sa"
    node             = "k8s-node-sa"
  }
}

resource "yandex_iam_service_account" "master_sa" {
  folder_id   = var.folder_id
  description = "Service account for k8s cluster. Master"
  name        = local.service_account.master
}

resource "yandex_iam_service_account" "node_sa" {
  folder_id   = var.folder_id
  description = "Service account for nodes in k8s cluster. Nodes"
  name        = local.service_account.node
}



resource "yandex_resourcemanager_folder_iam_member" "master_sa_roles" {

  folder_id = var.folder_id
  for_each  = toset([
    "k8s.clusters.agent", // This role is needed to node groups
    "k8s.tunnelClusters.agent", // This role is needed to manage network policies
    "vpc.publicAdmin", // This role is needed to manage public ip addresses
    "load-balancer.admin", // This role is needed to manage load balancers
    "logging.writer", // This role is needed to write logs to log group
  ])
  role       = each.value
  member     = "serviceAccount:${yandex_iam_service_account.master_sa.id}"
  depends_on = [
    yandex_iam_service_account.master_sa,
  ]
  sleep_after = 5
}

resource "yandex_resourcemanager_folder_iam_member" "node_sa_roles" {
  folder_id = var.folder_id
  for_each  = toset([
    "container-registry.images.puller", // This role is needed to pull images from k8s cluster registry
    "kms.keys.encrypterDecrypter", // This role is needed to decrypt on nodes secrets created by Cilium
  ])
  role       = each.value
  member     = "serviceAccount:${yandex_iam_service_account.node_sa.id}"
  depends_on = [
    yandex_iam_service_account.node_sa,
  ]
  sleep_after = 5
}

