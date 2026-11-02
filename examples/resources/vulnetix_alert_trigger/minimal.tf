# Only the arguments the kind requires. From config-schemas examples/AlertTrigger/minimal.yaml.

resource "vulnetix_alert_trigger" "dns_investigate" {
  name = "dns-investigate"
  destination = {
    kind = "workflow"
  }
  search_ref = "dynamic-dns-lookups"
}
