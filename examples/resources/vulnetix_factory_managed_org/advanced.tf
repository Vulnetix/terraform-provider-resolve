# Every argument the kind supports. From config-schemas examples/FactoryManagedOrg/advanced.yaml.

resource "vulnetix_factory_managed_org" "acme_platform" {
  name    = "acme-platform"
  account = "Acme-Platform"
  labels = {
    team = "platform"
  }
  annotations = {
    ticket = "https://tracker.example.com/PLAT-77"
  }
}

resource "vulnetix_factory_managed_org" "codeberg_org_acme" {
  name    = "codeberg-org-acme"
  account = "codeberg.org/acme"
}
