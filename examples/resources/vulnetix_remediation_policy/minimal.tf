# Only the arguments the kind requires. From config-schemas examples/RemediationPolicy/minimal.yaml.

resource "vulnetix_remediation_policy" "default" {
  standard = "nist-800-40"
}
