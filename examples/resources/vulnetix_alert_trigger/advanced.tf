# Every argument the kind supports. From config-schemas examples/AlertTrigger/advanced.yaml.

resource "vulnetix_alert_trigger" "dns_investigate" {
  name = "dns-investigate"
  destination = {
    kind     = "workflow"
    severity = "warning"
  }
  display_name = "Dynamic DNS lookups from one host"
  enabled      = true
  every_sec    = 300
  levers = {
    ask          = "always"
    columns      = ["ts", "src_ip", "qname"]
    cooldown_sec = 3600
    count_by     = "src_ip"
    hours        = "business"
    max_per_hour = 4
    min_rows     = 3
    rows_cap     = 25
    window_sec   = 900
    zone         = "Australia/Melbourne"
  }
  search_ref = "dynamic-dns-lookups"
  labels = {
    team = "secops"
  }
  annotations = {
    runbook = "https://wiki.example.com/secops/dynamic-dns"
  }
}

resource "vulnetix_alert_trigger" "root_login_page" {
  name = "root-login-page"
  destination = {
    connections = ["secops-slack"]
    kind        = "messaging"
    severity    = "critical"
  }
  display_name = "AWS root console login"
  every_sec    = 60
  levers = {
    cooldown_sec = 900
    hours        = "always"
    window_sec   = 300
  }
  rule_ref = "sigmahq-aws-root-login"
}

resource "vulnetix_alert_trigger" "failed_logins_digest" {
  name = "failed-logins-digest"
  destination = {
    emails   = ["secops@example.com", "oncall@example.com"]
    kind     = "email"
    severity = "notice"
  }
  display_name = "Failed SSO logins outside business hours"
  enabled      = false
  every_sec    = 3600
  levers = {
    count_by   = "user_name"
    hours      = "outside_business"
    min_rows   = 20
    rows_cap   = 50
    window_sec = 3600
    zone       = "Europe/London"
  }
  search_ref = "failed-sso-logins"
}
