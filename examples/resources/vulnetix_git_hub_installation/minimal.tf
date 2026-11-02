# Only the arguments the kind requires. From config-schemas examples/GitHubInstallation/minimal.yaml.

resource "vulnetix_git_hub_installation" "acme" {
  name            = "acme"
  installation_id = 51234567
}
