resource "kubernetes_namespace" "env" {
  metadata {
    labels = {
    }

    name = var.env
  }
}

// istio support
resource "kubectl_manifest" "istio_manifest_gateway" {
    yaml_body = templatefile("./gateway.yaml",
        {env = var.env}
    )
}

// file virtual_services.yaml includes many items to be created.
data "kubectl_path_documents" "istio_files_virtual_services" {
    pattern = "./virtual_services.yaml"
    vars = {env = var.env}
}

resource "kubectl_manifest" "istio_manifest_virtual_services" {
    count = length(data.kubectl_path_documents.istio_files_virtual_services.documents)
    yaml_body = element(data.kubectl_path_documents.istio_files_virtual_services.documents, count.index)
}

resource "kubernetes_service" "cart" {
    metadata {
        name = "cart"

        labels = {
            app = "cart"
        }

        annotations = {
        }

        namespace = var.env
    }

    spec {
        selector = {
            app = "cart"
        }

        port {
            port        = 1325
            target_port = 1325
        }

        type = "NodePort"
    }
}
/*
resource "kubernetes_service" "customer" {
    metadata {
        name = "customer"

        labels = {
            app = "customer"
        }

        annotations = {
            "field.cattle.io/publicEndpoints" = ""
        }

        namespace = var.env
    }

    spec {
        selector = {
            app = "customer"
        }

        port {
            port        = 5000
            target_port = 5000
        }

        type = "NodePort"
    }
}

resource "kubernetes_service" "product" {
    metadata {
        name = "product"

        labels = {
            app = "product"
        }

        annotations = {
            "field.cattle.io/publicEndpoints" = ""
        }

        namespace = var.env
    }

    spec {
        selector = {
            app = "product"
        }

        port {
            port        = 5000
            target_port = 5000
        }

        type = "NodePort"
    }
}*/

resource "kubernetes_deployment" "cart" {
    metadata {
        name = "cart"

        annotations = {
        }

        namespace = var.env
    }

    spec {
        progress_deadline_seconds = 6000

        replicas = 2

        selector {
            match_labels = {
                app = "cart"
            }
        }
    
        template {
            metadata {
                name = "cart"

                labels = {
                    app = "cart"
                }
            }
        
            spec {
                container {
                    name = "cart"
                    image = "${var.docker_registry}cart:${var.env}"
                    image_pull_policy = "Always"

                    liveness_probe {
                        http_get {
                            path = "/health"
                            port = "1325"
                        }

                        initial_delay_seconds = 1
                        period_seconds = 1
                    }

                    env {
                        name = "ENV"
                        value = var.env
                    }
                    env {
                        name = "STAGE"
                        value = var.stage
                    }
                }
            }
        }
    }

    timeouts {
        create = "10m"
        update = "10m"
        delete = "10m"
    }
}

/*
resource "kubernetes_deployment" "customer" {
    metadata {
        name = "customer"

        annotations = {
            "field.cattle.io/publicEndpoints" = ""
        }

        namespace = var.env
    }

    spec {
        progress_deadline_seconds = 6000

        replicas = 3

        selector {
            match_labels = {
                app = "customer"
            }
        }
    
        template {
            metadata {
                name = "customer"

                labels = {
                    app = "customer"
                }
            }
        
            spec {
                container {
                    name = "customer"
                    image = "${var.docker_registry}/service_customer:${var.env}"
                    image_pull_policy = "Always"

                    liveness_probe {
                        http_get {
                            path = "/health"
                            port = "5000"
                        }

                        initial_delay_seconds = 1
                        period_seconds = 1
                    }

                    env {
                        name = "ENV"
                        value = var.env
                    }
                    env {
                        name = "DB_HOST"
                        value = "127.0.0.1"
                    }
                }
            }
        }
    }

    lifecycle {
        ignore_changes = [
            metadata[0].annotations["field.cattle.io/publicEndpoints"]
        ]
    }

    timeouts {
        create = "1h"
        update = "1h"
        delete = "10m"
    }
}

resource "kubernetes_deployment" "product" {
    metadata {
        name = "product"

        annotations = {
            "field.cattle.io/publicEndpoints" = ""
        }

        namespace = var.env
    }

    spec {
        progress_deadline_seconds = 6000

        replicas = 3

        selector {
            match_labels = {
                app = "product"
            }
        }
    
        template {
            metadata {
                name = "product"

                labels = {
                    app = "product"
                }
            }
        
            spec {
                container {
                    name = "product"
                    image = "${var.docker_registry}/service_product:${var.env}"
                    image_pull_policy = "Always"

                    liveness_probe {
                        http_get {
                            path = "/health"
                            port = "5000"
                        }

                        initial_delay_seconds = 1
                        period_seconds = 1
                    }

                    env {
                        name = "ENV"
                        value = var.env
                    }
                    env {
                        name = "DB_HOST"
                        value = "127.0.0.1"
                    }
                }
            }
        }
    }

    lifecycle {
        ignore_changes = [
            metadata[0].annotations["field.cattle.io/publicEndpoints"]
        ]
    }

    timeouts {
        create = "1h"
        update = "1h"
        delete = "10m"
    }
}*/
