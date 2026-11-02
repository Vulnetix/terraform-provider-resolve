# Every argument the kind supports. From config-schemas examples/WizConnection/advanced.yaml.

variable "default_client_secret" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_wiz_connection" "default" {
  client_id                = "4ba3c9f2d7e14a6b8c0d1e2f3a4b5c6d"
  client_secret_wo         = var.default_client_secret
  client_secret_wo_version = 1
  tenant_datacenter        = "us17"
  labels = {
    owner = "cloud-security"
  }
  annotations = {
    note = "The service account needs the admin:audit permission."
  }
}
