# Only the arguments the kind requires. From config-schemas examples/RiskStrategy/minimal.yaml.

resource "vulnetix_risk_strategy" "catch_all" {
  name = "catch-all"
  rules = [
    {
      metric   = "REMAINING"
      operator = "ALWAYS"
    },
  ]
}
