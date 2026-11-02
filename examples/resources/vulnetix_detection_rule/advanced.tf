# Every argument the kind supports. From config-schemas examples/DetectionRule/advanced.yaml.

resource "vulnetix_detection_rule" "dynamic_dns" {
  name               = "dynamic-dns"
  body               = "title: Dynamic DNS lookup\nid: 6f3e0c52-0d4b-4c8e-9a51-2a6b9a0e7d10\nlevel: medium\nlogsource:\n  category: dns_query\ndetection:\n  selection:\n    QueryName|endswith:\n      - .duckdns.org\n      - .no-ip.org\n  condition: selection\n"
  display_name       = "Dynamic DNS lookup"
  enabled            = true
  interval_sec       = 120
  language           = "sigma"
  lateness_sec       = 300
  level              = "medium"
  max_alerts_per_run = 50
  severity_override  = "warning"
  suppress_sec       = 1800
  suppressions = [
    {
      field = "QueryName"
      until = "2026-12-31"
      value = "updates.duckdns.org"
    },
    {
      field = "src_ip"
      value = "10.20.0.15"
    },
  ]
  target = {
    connections = ["secops-slack"]
    emails      = ["secops@example.com"]
  }
  threshold = {
    by         = ["src_ip"]
    count      = 5
    window_sec = 600
  }
  labels = {
    mitre-tactic = "command-and-control"
    team         = "secops"
  }
  annotations = {
    runbook = "https://wiki.example.com/secops/dynamic-dns"
  }
}

resource "vulnetix_detection_rule" "brute_force_file" {
  name         = "brute-force-file"
  display_name = "SSH brute force"
  enabled      = false
  file         = "siem/rules/brute-force.yml"
  language     = "sigma"
  level        = "high"
}

resource "vulnetix_detection_rule" "dns_search_rule" {
  name         = "dns-search-rule"
  interval_sec = 600
  level        = "low"
  search_ref   = "dynamic-dns-lookups"
}

resource "vulnetix_detection_rule" "log4shell_countermeasure" {
  name              = "log4shell-countermeasure"
  countermeasure    = "3b1f8a62-7d2e-4c41-9a0b-5e6f7c8d9e01"
  interval_sec      = 300
  severity_override = "critical"
  target = {
    connections = ["secops-slack"]
  }
}

resource "vulnetix_detection_rule" "sigmahq_aws_root_login" {
  name         = "sigmahq-aws-root-login"
  level        = "high"
  sigmahq      = "rules/cloud/aws/cloudtrail/aws_root_account_usage.yml"
  suppress_sec = 3600
}

resource "vulnetix_detection_rule" "indicator_domains" {
  name               = "indicator-domains"
  display_name       = "Indicator domain resolved"
  ioc                = "ioc-dns"
  lateness_sec       = 120
  max_alerts_per_run = 100
}
