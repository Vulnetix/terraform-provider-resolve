# Only the arguments the kind requires. From config-schemas examples/MessagingConnection/minimal.yaml.

resource "vulnetix_messaging_connection" "security_slack" {
  name = "security-slack"
  type = "slack"
}
