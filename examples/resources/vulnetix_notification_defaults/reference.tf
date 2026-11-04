# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/NotificationDefaults/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_notification_defaults" "default" {
  push = {
  }
}
