# Every argument the kind supports. From config-schemas examples/QualityGate/advanced.yaml.

resource "vulnetix_quality_gate" "default" {
  block_malware  = true
  block_unpinned = true
  cooldown_days  = 7
  eol = {
    block = true
    severities = {
      next_quarter   = "skip"
      retired        = "critical"
      this_quarter   = "low"
      within_30_days = "high"
    }
  }
  exploits = "active"
  sca_autofix = {
    max_major_bump = 1
    strategy       = "stable"
  }
  severity    = "high"
  version_lag = 1
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2212"
  }
}
