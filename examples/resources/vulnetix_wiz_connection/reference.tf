# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/WizConnection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_wiz_connection" "default" {
  client_id         = "4ba3c9f2d7e14a6b8c0d1e2f3a4b5c6d"
  tenant_datacenter = "us17"
}
