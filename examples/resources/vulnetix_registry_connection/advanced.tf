# Every argument the kind supports. From config-schemas examples/RegistryConnection/advanced.yaml.

variable "acme_artifactory_npm_secret" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "acme_registry_secret" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_registry_connection" "acme_artifactory_npm" {
  name              = "acme-artifactory-npm"
  ecosystem         = "npm"
  enabled           = true
  scope             = ["@acme", "@acme-internal"]
  secret_wo         = var.acme_artifactory_npm_secret
  secret_wo_version = 1
  url               = "https://artifactory.acme.example/artifactory/api/npm/npm-local/"
  username          = "vulnetix-scanner"
  labels = {
    team = "platform"
  }
  annotations = {
    ticket = "https://tracker.example.com/PLAT-482"
  }
}

resource "vulnetix_registry_connection" "acme_registry" {
  name              = "acme-registry"
  ecosystem         = "oci"
  scope             = ["acme/payments", "acme/platform"]
  secret_wo         = var.acme_registry_secret
  secret_wo_version = 1
  url               = "registry.acme.example"
}
