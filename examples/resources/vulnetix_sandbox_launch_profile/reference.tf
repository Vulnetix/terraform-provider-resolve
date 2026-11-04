# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/SandboxLaunchProfile/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_sandbox_launch_profile" "dev_sandbox" {
  name    = "dev-sandbox"
  sandbox = "3c9e1f2a-7b4d-4e6f-8a1b-2c3d4e5f6a7b"
}
