# Only the arguments the kind requires. From config-schemas examples/ServiceAccount/minimal.yaml.

resource "vulnetix_service_account" "ci_scanner" {
  name = "ci-scanner"
}
