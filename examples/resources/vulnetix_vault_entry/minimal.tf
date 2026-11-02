# Only the arguments the kind requires. From config-schemas examples/VaultEntry/minimal.yaml.

resource "vulnetix_vault_entry" "npm_token" {
  name       = "npm-token"
  entry_name = "NPM_TOKEN"
  type       = "env"
}
