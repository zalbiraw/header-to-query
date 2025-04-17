terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
  }
}

variable "namespace" {
  description = "The namespace for the config map"
  type        = string
}

resource "kubernetes_config_map" "header_to_query" {
  metadata {
    name      = "header-to-query-plugin"
    namespace = var.namespace
  }

  data = merge({
    ".golangci.yml" = file("${path.module}/.golangci.yml")
    ".traefik.yml" = file("${path.module}/.traefik.yml")
    "go.mod" = file("${path.module}/go.mod")
    "headertoquery.go" = file("${path.module}/headertoquery.go")
    "Makefile" = file("${path.module}/Makefile")
    "headertoquery_test.go" = file("${path.module}/headertoquery_test.go")
  })
}
