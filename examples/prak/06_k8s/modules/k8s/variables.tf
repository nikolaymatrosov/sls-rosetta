variable "cloud_id" {
  description = "Cloud ID"
  type        = string
}
variable "folder_id" {
  description = "Folder ID"
  type        = string
}
variable "zone" {
  description = "Zone"
  type        = string
  default     = "ru-central1-a"
}
variable "master_version" {
  type        = string
  description = "Kubernetes version for master nodes"
  default     = "1.28"
}

variable "nodes_version" {
  type        = string
  description = "Kubernetes version for nodes"
  default     = "1.28"
}
variable "node_groups_defaults" {
  description = "Map of common default values for Node groups."
  type        = map(any)
  default     = {
    platform_id   = "standard-v3"
    node_cores    = 4
    node_memory   = 8
    node_gpus     = 0
    core_fraction = 100
    disk_type     = "network-ssd"
    disk_size     = 32
    preemptible   = false
    nat           = false
    ipv4          = true
    ipv6          = false
    size          = 1
  }
}

variable "network_name" {
  description = "Name of the network"
  type        = string
  default     = "k8s-network"
}