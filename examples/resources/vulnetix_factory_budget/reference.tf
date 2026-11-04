# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/FactoryBudget/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_factory_budget" "org_monthly_tokens" {
  name            = "org-monthly-tokens"
  action_at_limit = "alert_only"
  applies_to      = "org"
  limit           = 20000000
  unit            = "tokens"
}
