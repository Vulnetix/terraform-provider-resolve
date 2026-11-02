# Only the arguments the kind requires. From config-schemas examples/FactoryManagedOrg/minimal.yaml.

resource "vulnetix_factory_managed_org" "acme_platform" {
  name    = "acme-platform"
  account = "acme-platform"
}
