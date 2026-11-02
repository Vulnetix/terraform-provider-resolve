# Only the arguments the kind requires. From config-schemas examples/FactoryTrigger/minimal.yaml.

resource "vulnetix_factory_trigger" "run_now" {
  name     = "run-now"
  start_at = ["scan"]
  type     = "manual"
  workflow = "security-patch-train"
}
