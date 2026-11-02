# Only the arguments the kind requires. From config-schemas examples/AiGovernanceList/minimal.yaml.

resource "vulnetix_ai_governance_list" "deny_shell_tools" {
  name     = "deny-shell-tools"
  action   = "block"
  family   = "tool"
  list     = "deny"
  patterns = ["Bash"]
}
