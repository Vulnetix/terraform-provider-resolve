# Only the arguments the kind requires. From config-schemas examples/WizConnection/minimal.yaml.

resource "vulnetix_wiz_connection" "default" {
  client_id         = "4ba3c9f2d7e14a6b8c0d1e2f3a4b5c6d"
  tenant_datacenter = "us17"
}
