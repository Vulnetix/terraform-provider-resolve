# Every argument the kind supports. From config-schemas examples/JailExemption/advanced.yaml.

resource "vulnetix_jail_exemption" "payments_moment_migration" {
  name       = "payments-moment-migration"
  expires    = "2026-12-31T17:00:00Z"
  reason     = "Migration to date-fns is scheduled for the Q1 release train."
  repository = "acme/payments-api"
  rule       = "Off moment.js"
  labels = {
    team = "payments"
  }
  annotations = {
    ticket = "https://tracker.example.com/PAY-1432"
  }
}
