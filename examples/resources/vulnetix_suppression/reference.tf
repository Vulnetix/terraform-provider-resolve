# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/Suppression/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_suppression" "lodash_prototype_pollution" {
  name = "lodash-prototype-pollution"
  scope = {
    rule = "CVE-2021-23337"
  }
  type = "risk_accepted"
}
