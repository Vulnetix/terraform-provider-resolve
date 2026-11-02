# Every argument the kind supports. From config-schemas examples/GitHubInstallation/advanced.yaml.

resource "vulnetix_git_hub_installation" "acme" {
  name            = "acme"
  installation_id = 51234567
  labels = {
    owner = "platform"
  }
  annotations = {
    note = "Repository selection is managed on GitHub."
  }
}
