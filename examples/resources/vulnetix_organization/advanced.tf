# Every argument the kind supports. From config-schemas examples/Organization/advanced.yaml.

resource "vulnetix_organization" "default" {
  company_email_domains = ["acme.com", "acme.co.uk"]
  display_name          = "Acme Corporation"
  labels = {
    team = "platform"
  }
  annotations = {
    ticket = "https://tracker.example.com/ORG-12"
  }
}
