# Only the arguments the kind requires. From config-schemas examples/SandboxEgressPolicy/minimal.yaml.

resource "vulnetix_sandbox_egress_policy" "dev_sandbox_egress" {
  name    = "dev-sandbox-egress"
  sandbox = "3c9e1f2a-7b4d-4e6f-8a1b-2c3d4e5f6a7b"
}
