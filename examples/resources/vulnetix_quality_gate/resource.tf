# The tenant's quality gate. There is one per tenant, so name is always
# "default" and can be left out.
resource "vulnetix_quality_gate" "default" {
  severity       = "high"
  exploits       = "active"
  block_malware  = true
  block_unpinned = true
  cooldown_days  = 7

  eol = {
    block = true
    severities = {
      next_quarter   = "skip"
      this_quarter   = "low"
      within_30_days = "high"
      retired        = "critical"
    }
  }

  sca_autofix = {
    strategy       = "safest"
    max_major_bump = 0
  }
}
