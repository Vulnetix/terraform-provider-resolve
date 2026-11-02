# Only the arguments the kind requires. From config-schemas examples/FactorySchedule/minimal.yaml.

resource "vulnetix_factory_schedule" "nightly" {
  name     = "nightly"
  cron     = "0 3 * * *"
  start_at = ["scan"]
  workflow = "security-patch-train"
}
