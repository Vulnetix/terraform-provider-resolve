# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/RuntimeSensorPolicy/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_runtime_sensor_policy" "engineering_laptops" {
  name                    = "engineering-laptops"
  allow_major_update      = false
  allow_user_stop         = false
  auto_update             = true
  credential_mode         = "device-key"
  enabled                 = true
  heartbeat_alert_minutes = 15
  hold                    = "15m"
  inventory_every         = "6h"
  mode                    = "monitor"
  priority                = 100
  protection              = "hardened"
}
