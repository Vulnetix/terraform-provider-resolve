# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/QualityGate/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_quality_gate" "default" {
  eol = {
    severities = {
      next_quarter   = "low"
      retired        = "critical"
      this_quarter   = "medium"
      within_30_days = "high"
    }
  }
}
