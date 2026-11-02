# Every argument the kind supports. From config-schemas examples/PackageFirewallPolicy/advanced.yaml.

resource "vulnetix_package_firewall_policy" "default" {
  block = {
    active     = true
    bad_actors = true
    eol        = true
    kev        = true
    malware    = true
    poc        = false
    weaponized = true
  }
  cooldown_days = 3
  ecosystems = {
    go = {
      mode = "warn"
    }
    npm = {
      mode = "allowlist_only"
    }
    pypi = {
      mode = "enforce"
    }
  }
  mirrors = [
    {
      active    = true
      ecosystem = "go"
      priority  = 0
      url       = "https://proxy.golang.org"
    },
    {
      active    = false
      ecosystem = "go"
      priority  = 1
      url       = "https://goproxy.io"
    },
    {
      ecosystem = "npm"
      priority  = 0
      url       = "https://registry.npmjs.org"
    },
    {
      ecosystem = "pypi"
      priority  = 0
      url       = "https://pypi.org/simple"
    },
  ]
  mode = "enforce"
  review = {
    stale_decision_days = 180
    stale_usage_days    = 60
  }
  thresholds = {
    cess = 8
    cvss = 9
    epss = 0.5
  }
  version_lag = 10
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2210"
  }
}
