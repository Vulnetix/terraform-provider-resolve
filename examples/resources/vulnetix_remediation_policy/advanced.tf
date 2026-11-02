# Every argument the kind supports. From config-schemas examples/RemediationPolicy/advanced.yaml.

resource "vulnetix_remediation_policy" "default" {
  allowance_days = {
    critical = 2
    high     = 7
    low      = 30
    medium   = 14
  }
  description          = "Tighter than NIST for internet-facing estates."
  exposure_window_days = 14
  remediation_days = {
    critical = 7
    high     = 30
    low      = 180
    medium   = 90
  }
  severity_thresholds = {
    high   = 8.9
    low    = 4
    medium = 6.9
  }
  suppression_review_max_age_days = 180
  threat_window_days              = 30
  triage_threshold_days           = 2
  labels = {
    team = "platform-security"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2213"
  }
}

# Or, instead (one per tenant, so declare only one):
#
# resource "vulnetix_remediation_policy" "default_2" {
#   standard = "fedramp"
# }
