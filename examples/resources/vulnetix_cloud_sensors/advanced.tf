# Every argument the kind supports. From config-schemas examples/CloudSensors/advanced.yaml.

resource "vulnetix_cloud_sensors" "estate" {
  name    = "estate"
  account = "*"
  deploy  = true
  sources = ["inspector", "security-hub", "guardduty", "config", "macie", "access-analyzer", "waf", "compute", "compute-artifacts", "network-logs", "eks-workloads"]
  labels = {
    scope = "all-accounts"
  }
  annotations = {
    note = "Only account \"*\" applies today; a per-account document validates but is not stored yet."
  }
}
