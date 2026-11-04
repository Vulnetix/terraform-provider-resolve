# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/ScheduledReport/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_scheduled_report" "daily_dns_digest" {
  name = "daily-dns-digest"
  at   = "08:00"
  destination = {
    connections = ["secops-slack"]
    kind        = "messaging"
  }
  enabled = true
  format = {
    csv    = true
    layout = "table_chart"
    link   = true
    query  = false
    rows   = 20
  }
  monthday       = 1
  repeat         = "daily"
  rows_window_ms = 86400000
  search_ref     = "dynamic-dns-lookups"
  since_previous = false
  weekday        = 1
  when_empty     = "skip"
  zone           = "UTC"
}
