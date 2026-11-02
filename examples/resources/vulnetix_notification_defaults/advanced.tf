# Every argument the kind supports. From config-schemas examples/NotificationDefaults/advanced.yaml.

resource "vulnetix_notification_defaults" "default" {
  push = {
    "*"                     = false
    "factory.run.*"         = true
    "factory.run.succeeded" = false
    "tier:error"            = true
    "tier:needs_you"        = true
  }
  annotations = {
    ticket = "https://tracker.example.com/OPS-77"
  }
}
