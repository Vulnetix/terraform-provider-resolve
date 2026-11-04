# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/SsvcPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_ssvc_policy" "default" {
  automatable = "cvss >= 9.0 || epss >= 0.7"
  exploitation = {
    active = "kev || exploitMaturity in ['active', 'weaponized']"
    poc    = "exploitMaturity == 'poc' || exploitCount > 0"
  }
  mission_well_being = {
    critical     = "high"
    high         = "high"
    low          = "low"
    medium       = "medium"
    unclassified = "medium"
  }
  outcomes = {
    act = {
      due_days = 0
    }
    attend = {
      due_days = 0
    }
    track = {
      due_days = 0
    }
    track_star = {
      due_days = 0
    }
  }
  technical_impact = "cvss >= 9.0 || malicious"
}
