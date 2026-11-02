# Only the arguments the kind requires. From config-schemas examples/SavedSearch/minimal.yaml.

resource "vulnetix_saved_search" "dynamic_dns_lookups" {
  name = "dynamic-dns-lookups"
  text = "dns_events\n| where qname endswith ('.duckdns.org', '.no-ip.org', '.ddns.net')\n"
}
