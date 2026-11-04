# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/RepositorySettings/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_repository_settings" "acme_payments_api" {
  name = "acme-payments-api"
  asset_labels = {
  }
  autonomy           = "draft_pr"
  defaults_off       = []
  kanban_mode        = "off"
  protected_branches = true
  repository         = "acme/payments-api"
  tea_publish        = false
}
