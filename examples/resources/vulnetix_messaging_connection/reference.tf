# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/MessagingConnection/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_messaging_connection" "security_slack" {
  name   = "security-slack"
  paused = false
  type   = "slack"
}
