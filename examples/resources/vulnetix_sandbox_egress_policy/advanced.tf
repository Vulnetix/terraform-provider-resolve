# Every argument the kind supports. From config-schemas examples/SandboxEgressPolicy/advanced.yaml.

resource "vulnetix_sandbox_egress_policy" "dev_sandbox_egress" {
  name    = "dev-sandbox-egress"
  allow   = ["api.github.com", "*.githubusercontent.com", "registry.npmjs.org"]
  deny    = ["pastebin.com"]
  mode    = "deny_except_allow"
  sandbox = "3c9e1f2a-7b4d-4e6f-8a1b-2c3d4e5f6a7b"
  annotations = {
    ticket = "https://tracker.example.com/SEC-88"
  }
}
