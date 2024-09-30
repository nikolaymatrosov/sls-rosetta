resource "kubernetes_service" "nginx_service" {
  metadata {
    name      = "nginx-service"
    namespace = "praktikum"

    labels = {
      run = "my-nginx"
    }
  }

  spec {
    port {
      port = 80
    }

    selector = {
      run = "my-nginx"
    }
  }
}

