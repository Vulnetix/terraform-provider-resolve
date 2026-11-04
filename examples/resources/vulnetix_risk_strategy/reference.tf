# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/RiskStrategy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_risk_strategy" "catch_all" {
  name   = "catch-all"
  active = false
  rules = [
    {
      metric   = "REMAINING"
      operator = "ALWAYS"
    },
  ]
}
