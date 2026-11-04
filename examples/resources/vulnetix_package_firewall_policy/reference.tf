# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/PackageFirewallPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_package_firewall_policy" "default" {
  block = {
    active     = false
    bad_actors = false
    eol        = false
    kev        = false
    malware    = false
    poc        = false
    weaponized = false
  }
  cooldown_days = 0
  mode          = "enforce"
  review = {
    stale_decision_days = 90
    stale_usage_days    = 90
  }
  thresholds = {
    cess = 0
    cvss = 0
    epss = 0
  }
  version_lag = 0
}
