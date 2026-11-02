# Only the arguments the kind requires. From config-schemas examples/LogSource/minimal.yaml.

resource "vulnetix_log_source" "sensor_events" {
  name = "sensor-events"
  type = "events"
}
