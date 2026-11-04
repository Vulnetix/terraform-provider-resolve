# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/AlertTrigger/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_alert_trigger" "dns_investigate" {
  name = "dns-investigate"
  destination = {
    kind = "workflow"
  }
  enabled   = true
  every_sec = 300
  levers = {
    ask          = "workflow"
    cooldown_sec = 3600
    hours        = "always"
    max_per_hour = 4
    min_rows     = 1
    rows_cap     = 25
    window_sec   = 300
    zone         = "UTC"
  }
  search_ref = "dynamic-dns-lookups"
}
