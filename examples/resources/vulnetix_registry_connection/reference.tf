# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/RegistryConnection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_registry_connection" "github_packages" {
  name      = "github-packages"
  ecosystem = "npm"
  enabled   = true
  scope     = []
  url       = "https://npm.pkg.github.com/"
}
