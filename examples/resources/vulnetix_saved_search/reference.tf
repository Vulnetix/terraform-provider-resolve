# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/SavedSearch/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_saved_search" "dynamic_dns_lookups" {
  name     = "dynamic-dns-lookups"
  mode     = "dsl"
  range_ms = 86400000
  text     = "dns_events\n| where qname endswith ('.duckdns.org', '.no-ip.org', '.ddns.net')\n"
}
