# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/GitHubInstallation/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_git_hub_installation" "acme" {
  name            = "acme"
  installation_id = 51234567
}
