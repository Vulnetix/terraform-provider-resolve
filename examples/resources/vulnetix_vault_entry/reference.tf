# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/VaultEntry/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_vault_entry" "npm_token" {
  name       = "npm-token"
  entry_name = "NPM_TOKEN"
  type       = "env"
}
