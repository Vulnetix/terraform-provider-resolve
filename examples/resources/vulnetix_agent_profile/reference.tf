# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/AgentProfile/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_agent_profile" "reviewer" {
  name = "reviewer"
  body = "---\nid: 6f1c2b8a-3d4e-4f5a-9b6c-7d8e9f0a1b2c\nname: reviewer\n---\nYou review pull requests for correctness and security.\n"
}
