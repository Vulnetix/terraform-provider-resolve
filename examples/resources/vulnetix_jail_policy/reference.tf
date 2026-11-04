# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/JailPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_jail_policy" "org_default" {
  name                   = "org-default"
  default_staleness_days = 7
  enabled                = true
  enforcement_mode       = "enforce"
  on_stale               = "fail"
}
