# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/AiGovernanceList/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_ai_governance_list" "deny_shell_tools" {
  name     = "deny-shell-tools"
  action   = "block"
  enabled  = true
  family   = "tool"
  list     = "deny"
  patterns = ["Bash"]
  priority = 100
}
