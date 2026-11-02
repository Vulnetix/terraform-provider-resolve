# Only the arguments the kind requires. From config-schemas examples/RegistryConnection/minimal.yaml.

resource "vulnetix_registry_connection" "github_packages" {
  name      = "github-packages"
  ecosystem = "npm"
  url       = "https://npm.pkg.github.com/"
}
