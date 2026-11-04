# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/FactorySchedule/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_factory_schedule" "nightly" {
  name     = "nightly"
  cron     = "0 3 * * *"
  enabled  = true
  runs_on  = "aws"
  start_at = ["scan"]
  timezone = "UTC"
  workflow = "security-patch-train"
}
