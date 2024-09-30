resource "kubernetes_horizontal_pod_autoscaler_v2" "my_nginx_hpa" {
  metadata {
    name      = "my-nginx-hpa"
    namespace = "praktikum"
  }

  spec {
    scale_target_ref {
      kind        = "Deployment"
      name        = "my-nginx"
      api_version = "apps/v1"
    }

    min_replicas = 1
    max_replicas = 10

    metric {
      type = "Resource"

      resource {
        name = "cpu"

        target {
          type                = "Utilization"
          average_utilization = 50
        }
      }
    }
  }
}

