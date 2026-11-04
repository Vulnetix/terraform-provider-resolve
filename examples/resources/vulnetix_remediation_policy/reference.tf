# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/RemediationPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_remediation_policy" "default" {
  standard = "nist-800-40"
}
