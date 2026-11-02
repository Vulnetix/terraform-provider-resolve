# Only the arguments the kind requires. From config-schemas examples/ScheduledReport/minimal.yaml.

resource "vulnetix_scheduled_report" "daily_dns_digest" {
  name = "daily-dns-digest"
  destination = {
    connections = ["secops-slack"]
    kind        = "messaging"
  }
  search_ref = "dynamic-dns-lookups"
}
