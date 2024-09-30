module "k8s" {
  source = "./modules/k8s"

  cloud_id   = var.cloud_id
  folder_id  = var.folder_id
}

resource "kubernetes_namespace_v1" "praktikum" {
  metadata {
    name = "praktikum"
  }
}

resource "kubernetes_deployment_v1" "praktikum-nginx" {
  metadata {
    name = "my-nginx"
    labels = {
      test = "my-nginx"
    }
    namespace = kubernetes_namespace_v1.praktikum.metadata.0.name
  }

  spec {
    replicas = 3

    selector {
      match_labels = {
        test = "my-nginx"
      }
    }

    template {
      metadata {
        labels = {
          test = "my-nginx"
        }
      }

      spec {
        container {
          image = local.image
          name  = "nginx"

          resources {
            limits = {
              cpu    = "1"
              memory = "500Mi"
            }
            requests = {
              cpu    = "500m"
              memory = "256Mi"
            }
          }

          liveness_probe {
            http_get {
              path = "/"
              port = 80

            }

            initial_delay_seconds = 3
            period_seconds        = 3
          }
        }
      }
    }
  }
}
