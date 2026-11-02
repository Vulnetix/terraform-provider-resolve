# Every argument the kind supports. From config-schemas examples/FactoryBudget/advanced.yaml.

resource "vulnetix_factory_budget" "patch_train_daily" {
  name            = "patch-train-daily"
  action_at_limit = "pause_run"
  alert_at        = [50, 80, 100]
  applies_to      = "workflow"
  limit           = 1500000
  model_providers = "anthropic,openai"
  unit            = "tokens"
  window          = "day"
  workflow        = "security-patch-train"
  labels = {
    team = "appsec"
  }
}

resource "vulnetix_factory_budget" "patch_step" {
  name            = "patch-step"
  action_at_limit = "downgrade_to_fast"
  applies_to      = "node"
  limit           = 300000
  node            = "patch"
  window          = "node"
  workflow        = "security-patch-train"
}

resource "vulnetix_factory_budget" "payments_repo" {
  name            = "payments-repo"
  action_at_limit = "escalate_human"
  applies_to      = "repo"
  limit           = 5000
  repository      = "acme/payments"
  unit            = "calls"
  window          = "month"
}

resource "vulnetix_factory_budget" "acme_org" {
  name            = "acme-org"
  action_at_limit = "alert_only"
  applies_to      = "org"
  installation_id = 41234567
  limit           = 10000000
}
