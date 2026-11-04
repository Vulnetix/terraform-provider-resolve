# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/FactoryManagedOrg/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_factory_managed_org" "acme_platform" {
  name    = "acme-platform"
  account = "acme-platform"
}
