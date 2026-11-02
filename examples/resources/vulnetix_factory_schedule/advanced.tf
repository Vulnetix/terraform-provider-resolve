# Every argument the kind supports. From config-schemas examples/FactorySchedule/advanced.yaml.

resource "vulnetix_factory_schedule" "weekday_sweep" {
  name     = "weekday-sweep"
  cron     = "0 3 * * 1-5"
  enabled  = true
  overlap  = "skip_if_running"
  priority = 3
  runs_on  = "aws"
  start_at = ["scan"]
  timezone = "Australia/Melbourne"
  title    = "Weekday dependency sweep"
  workflow = "security-patch-train"
  annotations = {
    ticket = "https://tracker.acme.com/SEC-142"
  }
}
