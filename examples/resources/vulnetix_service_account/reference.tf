# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/ServiceAccount/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_service_account" "ci_scanner" {
  name = "ci-scanner"
  role = "member"
}
