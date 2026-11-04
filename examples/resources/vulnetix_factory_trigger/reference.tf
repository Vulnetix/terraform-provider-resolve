# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/FactoryTrigger/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_factory_trigger" "run_now" {
  name     = "run-now"
  start_at = ["scan"]
  type     = "manual"
  workflow = "security-patch-train"
}
