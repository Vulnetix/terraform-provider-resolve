# Only the arguments the kind requires. From config-schemas examples/FactoryBudget/minimal.yaml.

resource "vulnetix_factory_budget" "org_monthly_tokens" {
  name       = "org-monthly-tokens"
  applies_to = "org"
  limit      = 20000000
}
