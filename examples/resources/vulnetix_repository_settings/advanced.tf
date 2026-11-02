# Every argument the kind supports. From config-schemas examples/RepositorySettings/advanced.yaml.

resource "vulnetix_repository_settings" "acme_payments_api" {
  name = "acme-payments-api"
  asset_labels = {
    data = "pii"
    pci  = "in-scope"
    tier = "prod"
  }
  autonomy          = "suggest"
  criticality       = "critical"
  defaults_off      = ["lint-sweep", "test-sweep"]
  kanban_mode       = "skip_review"
  monthly_token_cap = 2000000
  owner = {
    kind = "team"
    name = "payments"
  }
  protected_branches = true
  repository         = "acme/payments-api"
  tea_publish        = true
  labels = {
    team = "payments"
  }
  annotations = {
    ticket = "https://tracker.example.com/PAY-1501"
  }
}
