# Only the arguments the kind requires. From config-schemas examples/RepositorySettings/minimal.yaml.

resource "vulnetix_repository_settings" "acme_payments_api" {
  name       = "acme-payments-api"
  autonomy   = "draft_pr"
  repository = "acme/payments-api"
}
