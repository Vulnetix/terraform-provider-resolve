# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/ComplianceTarget/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_compliance_target" "quarterly_owasp" {
  name      = "quarterly-owasp"
  cadence   = "quarterly"
  enabled   = true
  framework = "owasp_top_10"
  scope = {
    repositories = []
  }
}
