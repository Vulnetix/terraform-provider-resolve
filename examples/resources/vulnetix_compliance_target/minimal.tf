# Only the arguments the kind requires. From config-schemas examples/ComplianceTarget/minimal.yaml.

resource "vulnetix_compliance_target" "quarterly_owasp" {
  name      = "quarterly-owasp"
  framework = "owasp_top_10"
}
