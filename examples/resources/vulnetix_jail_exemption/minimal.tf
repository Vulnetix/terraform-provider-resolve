# Only the arguments the kind requires. From config-schemas examples/JailExemption/minimal.yaml.

resource "vulnetix_jail_exemption" "payments_moment_migration" {
  name       = "payments-moment-migration"
  reason     = "Migration to date-fns is scheduled for the Q1 release train."
  repository = "acme/payments-api"
}
