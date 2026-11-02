# Every argument the kind supports. From config-schemas examples/SavedSearch/advanced.yaml.

resource "vulnetix_saved_search" "dynamic_dns_lookups" {
  name         = "dynamic-dns-lookups"
  display_name = "Dynamic DNS lookups"
  mode         = "dsl"
  range_ms     = 86400000
  text         = "dns_events\n| where qname endswith ('.duckdns.org', '.no-ip.org', '.ddns.net')\n"
  labels = {
    team = "secops"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2215"
  }
}

resource "vulnetix_saved_search" "console_sign_ins" {
  name         = "console-sign-ins"
  display_name = "AWS console sign-ins"
  mode         = "sigma"
  range_ms     = 604800000
  text         = "logsource:\n  product: aws\n  service: cloudtrail\ndetection:\n  selection:\n    eventName: ConsoleLogin\n  condition: selection\n"
}
