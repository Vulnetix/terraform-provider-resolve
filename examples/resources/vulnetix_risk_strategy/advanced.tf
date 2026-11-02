# Every argument the kind supports. From config-schemas examples/RiskStrategy/advanced.yaml.

resource "vulnetix_risk_strategy" "exploit_first" {
  name        = "exploit-first"
  active      = true
  description = "Known exploitation first, then exposure, then score."
  rules = [
    {
      label    = "Malware"
      metric   = "MALWARE"
      operator = "IS_TRUE"
    },
    {
      label    = "CISA KEV"
      metric   = "KEV_CISA"
      operator = "IS_TRUE"
    },
    {
      enabled  = true
      metric   = "EXPLOIT_WEAPONIZED"
      operator = "IS_TRUE"
    },
    {
      label    = "Exposed package"
      metric   = "INTERNET_EXPOSED_PACKAGE"
      operator = "IS_TRUE"
    },
    {
      attack_path_id = "0b9d5f9e-3c51-4a0f-9d43-6f1f5f0c2a11"
      label          = "Edge to database"
      metric         = "ATTACK_PATH"
      operator       = "IS_TRUE"
    },
    {
      label  = "Reachable and likely"
      metric = "CEL"
      when   = "reachable && (epss >= 0.2 || exploitMaturity == 'active')"
    },
    {
      metric        = "EPSS"
      operator      = "GTE"
      threshold_low = 0.1
    },
    {
      metric         = "CVSS"
      operator       = "BETWEEN"
      threshold_high = 10
      threshold_low  = 7
    },
    {
      enabled  = false
      label    = "Public PoC"
      metric   = "EXPLOIT_POC"
      operator = "IS_TRUE"
    },
    {
      label    = "Everything else"
      metric   = "REMAINING"
      operator = "ALWAYS"
    },
  ]
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2214"
  }
}
