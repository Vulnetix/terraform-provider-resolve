# Every argument the kind supports. From config-schemas examples/SsvcPolicy/advanced.yaml.

resource "vulnetix_ssvc_policy" "default" {
  automatable = "cvss >= 9.0 || epss >= 0.5"
  exploitation = {
    active = "kev || exploitMaturity in ['active', 'weaponized']"
    poc    = "exploitMaturity == 'poc' || exploitCount > 0"
  }
  mapping = [
    {
      automatable        = ["no", "yes"]
      exploitation       = ["active"]
      mission_well_being = ["medium", "high"]
      outcome            = "Act"
      technical_impact   = ["total"]
    },
  ]
  mission_well_being = {
    critical     = "high"
    high         = "high"
    low          = "low"
    medium       = "medium"
    unclassified = "medium"
  }
  outcomes = {
    act = {
      due_days = 3
      priority = "immediate"
    }
    attend = {
      due_days = 14
      priority = "high"
    }
    track = {
      due_days = 90
      priority = "low"
    }
    track_star = {
      due_days = 60
      priority = "medium"
    }
  }
  overrides = [
    {
      name    = "card data in scope"
      outcome = "Act"
      when    = "labels['pci'] == 'in-scope' && severity in ['critical', 'high']"
    },
    {
      name    = "no fix and not reachable"
      outcome = "Track"
      when    = "!reachable && fixAvailability == 'no_fix'"
    },
  ]
  technical_impact = "cvss >= 9.0 || malicious"
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-3120"
  }
}
