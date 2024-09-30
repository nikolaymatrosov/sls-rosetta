output "cluster_id" {
  description = "Kubernetes cluster ID."
  value       = try(yandex_kubernetes_cluster.k8s_cluster.id, null)
}

output "cluster_name" {
  description = "Kubernetes cluster name."
  value       = try(yandex_kubernetes_cluster.k8s_cluster.name, null)
}

# public ip with kube config download command
output "external_cluster_cmd" {
  description = <<EOF
    Kubernetes cluster public IP address.
    Use the following command to download kube config and start working with Yandex Managed Kubernetes cluster:
    `$ yc managed-kubernetes cluster get-credentials --id <cluster_id> --external`
    This command will automatically add kube config for your user; after that, you will be able to test it with the
    `kubectl get cluster-info` command.
  EOF
  value       = "yc managed-kubernetes cluster get-credentials --id ${yandex_kubernetes_cluster.k8s_cluster.id} --external"
}


# The following output variables are for using with terraform providers such as kubernetes, helm, et cetera

# public cluster ip
output "external_v4_address" {
  description = <<EOF
    Kubernetes cluster external IP address.
  EOF
  value       = yandex_kubernetes_cluster.k8s_cluster.master[0].external_v4_address
}

# public cluster url
output "external_v4_endpoint" {
  description = <<EOF
    Kubernetes cluster external URL.
  EOF
  value       = yandex_kubernetes_cluster.k8s_cluster.master[0].external_v4_endpoint
}

# private cluster url
output "internal_v4_endpoint" {
  description = <<EOF
    Kubernetes cluster internal URL.
    Note: Kubernetes internal cluster nodes are available from the virtual machines in the same VPC as cluster nodes.
  EOF
  value       = yandex_kubernetes_cluster.k8s_cluster.master[0].internal_v4_endpoint
}

# cluster CA certificate
output "cluster_ca_certificate" {
  description = <<EOF
    Kubernetes cluster certificate.
  EOF
  value       = yandex_kubernetes_cluster.k8s_cluster.master[0].cluster_ca_certificate
}

output "network_id" {
  value = yandex_vpc_network.default.id
}

output "subnet_id" {
  value = {
    for k, v in yandex_vpc_subnet.default : k => v.id
  }
}

output "subnet_cidr" {
  value = {
    for k, v in yandex_vpc_subnet.default : k => v.v4_cidr_blocks
  }
}

output "node_group_id" {
  value = yandex_kubernetes_node_group.django_node_group.id
}

