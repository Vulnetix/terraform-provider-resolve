# Only the arguments the kind requires. From config-schemas examples/CloudSensors/minimal.yaml.

resource "vulnetix_cloud_sensors" "estate" {
  name    = "estate"
  sources = ["inspector", "security-hub"]
}
