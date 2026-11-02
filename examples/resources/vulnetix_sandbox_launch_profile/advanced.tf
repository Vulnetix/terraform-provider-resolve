# Every argument the kind supports. From config-schemas examples/SandboxLaunchProfile/advanced.yaml.

resource "vulnetix_sandbox_launch_profile" "dev_sandbox" {
  name       = "dev-sandbox"
  classifier = "jev"
  extensions = [
    {
      checksum = "sha256-AbCdEf0123456789AbCdEf0123456789AbCdEf01234="
      name     = "ripgrep"
      registry = "nixpkgs"
      version  = "14.1.0"
    },
  ]
  image   = "ghcr.io/vulnetix/pix-sandbox@sha256:2f4c6a8e0b1d3f5a7c9e1b3d5f7a9c1e3b5d7f9a1c3e5b7d9f1a3c5e7b9d1f3a"
  notes   = "Pinned to the October image until the Python 3.13 packs land."
  packs   = ["python", "node"]
  sandbox = "3c9e1f2a-7b4d-4e6f-8a1b-2c3d4e5f6a7b"
  version_policy = {
    pins = {
      ruff = "0.8.4"
    }
    strategy = "pinned"
  }
  web_controls         = true
  web_guardrails_off   = false
  web_project_settings = true
  labels = {
    team = "platform"
  }
}
