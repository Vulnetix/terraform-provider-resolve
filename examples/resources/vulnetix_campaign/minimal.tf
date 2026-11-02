# Only the arguments the kind requires. From config-schemas examples/Campaign/minimal.yaml.

resource "vulnetix_campaign" "log4shell" {
  name         = "log4shell"
  display_name = "Log4Shell clean-up"
  done = {
    criteria = [
      {
        type = "gone-from-scans"
      },
    ]
  }
  target = "vulnerabilities"
}
