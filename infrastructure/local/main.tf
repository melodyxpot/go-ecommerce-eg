terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "~> 1.13.3"
    }
  }
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
}
