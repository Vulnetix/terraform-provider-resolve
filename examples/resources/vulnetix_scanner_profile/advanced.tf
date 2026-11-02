# Every argument the kind supports. From config-schemas examples/ScannerProfile/advanced.yaml.

resource "vulnetix_scanner_profile" "default" {
  custom   = ["sca", "malscan", "secrets", "sast"]
  evidence = "sbom_only"
  scanners = "custom"
  slsa_min = "L2"
  labels = {
    team = "security"
  }
  annotations = {
    ticket = "https://tracker.example.com/SEC-210"
  }
}
