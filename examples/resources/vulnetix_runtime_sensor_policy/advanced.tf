# Every argument the kind supports. From config-schemas examples/RuntimeSensorPolicy/advanced.yaml.

resource "vulnetix_runtime_sensor_policy" "engineering_laptops" {
  name               = "engineering-laptops"
  allow_major_update = true
  allow_user_stop    = false
  credential_mode    = "device-key"
  description        = "Every engineer's laptop, in guard mode, plus the security team by address."
  dns_listen         = "127.0.0.1:5353"
  egress_listen      = "127.0.0.1:8181"
  enabled            = true
  feed_key           = "KymtMlKOJvj2IkhUbMpXRl/JBmmOwu/LIj8BekwmQXk="
  feeds = jsonencode([{
    feedRef = "edge-waf"
    }, "https://vulnetix.com/feeds/0b6f2c1e-4d3a-4e5f-8a9b-7c6d5e4f3a2b/9c8b7a6f-5e4d-4c3b-8a2f-1e0d9c8b7a6f", {
    sensorFeed = "7d6c5b4a-3f2e-4d1c-9b8a-7f6e5d4c3b2a"
  }])
  forward                 = ["siem"]
  heartbeat_alert_minutes = 30
  hold                    = "15m"
  inventory_every         = "6h"
  match = {
    domains          = ["acme.example"]
    emails           = ["alice@acme.example", "bob@acme.example"]
    mdm_integrations = ["5e4d3c2b-1a09-4f8e-9d7c-6b5a4f3e2d1c"]
  }
  mode       = "guard"
  priority   = 200
  protection = "hardened"
  labels = {
    fleet = "laptops"
  }
  annotations = {
    ticket = "https://jira.example.com/browse/SEC-2210"
  }
}

resource "vulnetix_runtime_sensor_policy" "servers_locked" {
  name            = "servers-locked"
  auto_update     = false
  credential_mode = "org-key"
  description     = "Production servers enrolled by MDM, blocking every match on a block-mode rule."
  feeds = jsonencode([{
    feedRef = "edge-suricata"
  }])
  heartbeat_alert_minutes = 5
  match = {
    mdm_integrations = ["5e4d3c2b-1a09-4f8e-9d7c-6b5a4f3e2d1c"]
  }
  mode       = "enforce"
  priority   = 10
  protection = "locked"
}
