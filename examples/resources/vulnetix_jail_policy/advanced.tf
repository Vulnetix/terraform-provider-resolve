# Every argument the kind supports. From config-schemas examples/JailPolicy/advanced.yaml.

resource "vulnetix_jail_policy" "payments_api" {
  name                   = "payments-api"
  default_staleness_days = 7
  description            = "Payments ships with no overdue criticals, no live secrets, and moment.js on its way out."
  enabled                = true
  enforcement_mode       = "enforce"
  on_stale               = "fail"
  repository             = "acme/payments-api"
  rules = [
    {
      aggregate      = "COUNT_PAST_SLA"
      category       = "SCA"
      ecosystem      = "npm"
      enabled        = true
      flags          = ["reachable", "hasFix"]
      kind           = "VULN"
      label          = "No reachable critical npm finding past its SLA"
      on_breach      = "BLOCK"
      on_stale       = "warn"
      operator       = "GT"
      ratchet        = false
      severity_floor = "critical"
      staleness_days = 3
      threshold_low  = 0
      tool_name      = "vulnetix"
      when           = "fixAvailability == 'available' && (kev || epss >= 0.1)"
    },
    {
      aggregate      = "COUNT"
      deadline       = "2027-03-31"
      kind           = "GOAL"
      label          = "Off moment.js"
      on_breach      = "WARN"
      operator       = "GT"
      ratchet        = true
      target_pattern = "pkg:npm/moment@*"
      threshold_low  = 0
    },
    {
      aggregate      = "COUNT"
      category       = "SECRETS"
      kind           = "HYGIENE"
      label          = "Secrets stay between 0 and 2"
      on_stale       = "pass"
      operator       = "OUTSIDE"
      threshold_high = 2
      threshold_low  = 0
    },
  ]
  labels = {
    team = "payments"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-1980"
  }
}
