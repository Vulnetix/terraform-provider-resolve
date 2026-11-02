# Every argument the kind supports. From config-schemas examples/ScheduledReport/advanced.yaml.

resource "vulnetix_scheduled_report" "weekly_console_sign_ins" {
  name = "weekly-console-sign-ins"
  at   = "08:00"
  destination = {
    emails  = ["secops@acme.example", "ciso@acme.example"]
    kind    = "email"
    subject = "Console sign-ins this week"
  }
  display_name = "Weekly console sign-ins"
  enabled      = true
  format = {
    csv    = true
    layout = "table_chart"
    link   = true
    query  = true
    rows   = 20
  }
  repeat         = "weekly"
  rows_window_ms = 604800000
  search_ref     = "console-sign-ins"
  weekday        = 1
  when_empty     = "all_clear"
  zone           = "Australia/Melbourne"
  labels = {
    team = "secops"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2216"
  }
}

resource "vulnetix_scheduled_report" "monthly_dns_digest" {
  name = "monthly-dns-digest"
  at   = "09:30"
  destination = {
    connections = ["secops-slack", "soc-teams"]
    kind        = "messaging"
    subject     = "Dynamic DNS lookups this month"
  }
  enabled = true
  format = {
    layout = "summary"
    link   = true
  }
  monthday       = 1
  repeat         = "monthly"
  search_ref     = "dynamic-dns-lookups"
  since_previous = true
  when_empty     = "skip"
  zone           = "UTC"
}
