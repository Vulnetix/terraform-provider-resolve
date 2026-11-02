# Only the arguments the kind requires. From config-schemas examples/NotificationDefaults/minimal.yaml.

resource "vulnetix_notification_defaults" "default" {
  push = {
  }
}
