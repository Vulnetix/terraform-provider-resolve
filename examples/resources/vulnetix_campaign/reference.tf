# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/Campaign/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_campaign" "log4shell" {
  name         = "log4shell"
  display_name = "Log4Shell clean-up"
  done = {
    criteria = [
      {
        type = "gone-from-scans"
      },
    ]
  }
  target = "vulnerabilities"
  timing = "frozen"
}
