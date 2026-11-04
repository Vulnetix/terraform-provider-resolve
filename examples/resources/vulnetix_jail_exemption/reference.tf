# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/JailExemption/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_jail_exemption" "payments_moment_migration" {
  name       = "payments-moment-migration"
  reason     = "Migration to date-fns is scheduled for the Q1 release train."
  repository = "acme/payments-api"
}
