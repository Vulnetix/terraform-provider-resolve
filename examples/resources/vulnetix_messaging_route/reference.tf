# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/MessagingRoute/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_messaging_route" "broken_builds" {
  name = "broken-builds"
  actor_filter = {
    mode = "any"
  }
  allow_sensitive_webhook = false
  digest = {
    mode = "immediate"
  }
  enabled     = true
  event_types = ["scan.build_broken"]
  filters = {
    findings = {
      ignore_test     = false
      new_only        = false
      only_direct     = false
      only_production = false
    }
  }
  include_pii  = false
  max_per_hour = 60
  min_severity = "info"
  mode         = "live"
}
